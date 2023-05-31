package job

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/reugn/go-quartz/quartz"
	"job/internal/tools/datex"
	"job/model"

	"time"
)

const (
	Wait    = "wait"
	Readly  = "readly"
	Running = "running"
	Stop    = "stop"
	Enable  = "enable"
	Reject  = "reject"

	Action     = "action"
	Compensate = "compensate"
)

var (
	JobScheduleNotFindError = errors.New("job schedule not exist")
	JobPlugNotFindError     = errors.New("job plug not exist")
	JobCronEmptyError       = errors.New("job cron not empty")
)

type DB struct {
	JobRecord   model.JobRecordModel
	JobSchedule model.JobScheduleModel
}

type Result struct {
	Data   any    `json:"data"`
	Key    string `json:"key"`
	Params any    `json:"params"`
	Error  string `json:"error"`
}

type node struct {
	plug  Plug   // 用户实现的插件逻辑
	key   int    // DbJob的key, 当dbJob在内存中创建时才会有key
	state string // 状态 wait  | readly | running
}

type Scheduler struct {
	// job的ctx
	ctx context.Context
	// db ctx 避免
	ctxDb  context.Context
	cancel context.CancelFunc
	quartz quartz.Scheduler
	db     *DB
	// 保存db-plug-key三者间的联系，核心
	nodeStore map[string]node
}

func NewScheduler(db *DB) *Scheduler {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &Scheduler{
		cancel:    cancelFunc,
		ctx:       ctx,
		ctxDb:     context.Background(),
		quartz:    quartz.NewStdScheduler(),
		db:        db,
		nodeStore: map[string]node{},
	}
}

func (s *Scheduler) Start() *Scheduler {
	s.quartz.Start(s.ctx)
	return s
}

func (s *Scheduler) Restart() {
	s.Stop()
	s.Start()
}

func (s *Scheduler) Plugs(plugs ...Plug) (errs []error, ok bool) {
	ctx := context.Background()
	errs = make([]error, 0)

	for _, plug := range plugs {
		key := plug.Key()
		if key == "" {
			errs = append(errs, errors.New("plug key error: key not empty"))
			continue
		}
		if _, exist := s.nodeStore[key]; exist {
			errs = append(errs, errors.New(fmt.Sprintf("plug key error: '%s' already exists", key)))
			continue
		}
		s.nodeStore[key] = node{
			plug:  plug,
			key:   0,
			state: Wait,
		}
	}
	if len(errs) > 0 {
		return errs, false
	}

	ok = true
	for _, plug := range plugs {
		err := s.AddScheduleDbJob(ctx, plug.Key())
		if err != nil {
			ok = false
			errs = append(errs, errors.Wrapf(err, "[key: %s]", plug.Key()))
		}
	}

	return errs, ok
}

func (s *Scheduler) IsStarted() bool {
	return s.quartz.IsStarted()
}

// AddScheduleDbJob 添加调度任务，id是任务key，即插件的key
func (s *Scheduler) AddScheduleDbJob(ctx context.Context, k string) error {
	jobInfo, err := s.db.JobSchedule.FindByKey(s.ctxDb, k)
	if jobInfo == nil {
		return JobScheduleNotFindError
	}
	if jobInfo.JobStatus != Enable {
		//_ = s.deleteJob(jobInfo.Key)
		return nil
	}

	cronTrigger, err := quartz.NewCronTriggerWithLoc(jobInfo.Cron, datex.TimeLocation)
	if err != nil {
		return JobCronEmptyError
	}

	job, err := s.warpJob(jobInfo)
	if err != nil {
		return err
	}

	return s.quartz.ScheduleJob(ctx, job, cronTrigger)
}

// StopJob key:唯一标识, 从数据库中和内存调度中停止Job, JobStatus=stop
func (s *Scheduler) StopJob(ctx context.Context, key string) error {
	err := s.db.JobSchedule.Update(s.ctxDb, &model.JobSchedule{
		Key:       key,
		JobStatus: Stop,
	})
	if err != nil {
		return err
	}
	err, _ = s.deleteJob(key)
	return err
}

func (s *Scheduler) RemoveJob(key string) (error, bool) {
	return s.deleteJob(key)
}

func (s *Scheduler) Clear(ctx context.Context) error {
	err := s.db.JobSchedule.Update(s.ctxDb, &model.JobSchedule{JobStatus: Stop})
	if err != nil {
		return err
	}
	s.quartz.Clear()
	return nil
}

func (s *Scheduler) Wait(ctx context.Context) {
	s.quartz.Wait(ctx)
}

// ForceStop 强制停止, 停止所有job, 正在执行的job会立即中止，可能会产生意料之外的情况
func (s *Scheduler) ForceStop() {
	s.cancel()
	s.Stop()
}

// Stop 停止所有job, 不会中止正在执行的job；优雅stop
func (s *Scheduler) Stop() {
	s.quartz.Stop()
}

// RestartJob 重置job
func (s *Scheduler) RestartJob(ctx context.Context, k string) error {
	err, _ := s.deleteJob(k)
	if err != nil {
		return err
	}
	return s.AddScheduleDbJob(ctx, k)
}

// ForceResetJob ResetJob 强制重置job, lock=1时会强制解锁
func (s *Scheduler) ForceResetJob(ctx context.Context, k string) error {
	jobInfo, _ := s.db.JobSchedule.FindByKey(s.ctxDb, k)
	if jobInfo == nil {
		return JobScheduleNotFindError
	}
	if jobInfo.IsLock() && !s.db.JobSchedule.UnLockByKey(s.ctxDb, k) {
		return errors.New("Force Restart job fail, unlock fail!")
	}

	return s.RestartJob(ctx, k)
}

// AddCompensate 执行补偿任务;params执行参数;background是否后台运行，true后台运行非阻塞， false阻塞
func (s *Scheduler) AddCompensate(ctx context.Context, key string, params map[string]any, background bool) (res any, err error) {
	jobInfo, _ := s.db.JobSchedule.FindByKey(s.ctxDb, key)
	if jobInfo == nil {
		return nil, JobScheduleNotFindError
	}

	plug := s.getPlugWithId(key)
	if plug == nil {
		return nil, JobPlugNotFindError
	}

	execFn := func() {
		startTime := time.Now()
		jobRecord := &model.JobRecord{
			JobId:     jobInfo.Id,
			StartTime: &startTime,
			EndTime:   nil,
			Result:    "",
			Status:    Running,
			UseMilli:  0,
			ExecType:  Compensate,
		}
		_ = s.db.JobRecord.Insert(s.ctxDb, jobRecord)

		func() {
			defer func() {
				if r := recover(); r != nil {
					err = errors.New(fmt.Sprint(r))
				}
			}()
			res, err = plug.Compensate(ctx, params)
		}()

		resStatus := Ok
		if err != nil {
			resStatus = Error
		}
		resultJson := NewResult(res, key, err, params).ToJson()
		endTime := time.Now()
		jobRecord.EndTime = &endTime
		jobRecord.Result = resultJson
		jobRecord.UseMilli = endTime.UnixMilli() - startTime.UnixMilli()
		jobRecord.Status = string(resStatus)
		_ = s.db.JobRecord.Update(s.ctxDb, jobRecord)
	}

	if background {
		go execFn()
	} else {
		execFn()
	}

	return res, err
}

func (s *Scheduler) getPlugWithId(k string) Plug {
	if n, ok := s.nodeStore[k]; ok {
		return n.plug
	}
	return nil
}

func (s *Scheduler) getJobKeyWithId(k string) (key int, exist bool) {
	if n, ok := s.nodeStore[k]; ok {
		return n.key, true
	}
	return 0, false
}

func (s *Scheduler) setNodeKey(k string, key int) {
	if n, ok := s.nodeStore[k]; ok {
		n.key = key
		s.nodeStore[k] = n
	}
}

func (s *Scheduler) setNodeState(k string, state string) {
	if n, ok := s.nodeStore[k]; ok {
		n.state = state
		s.nodeStore[k] = n
	}
}

func (s *Scheduler) deleteJob(k string) (err error, ok bool) {
	key, exist := s.getJobKeyWithId(k)
	if exist {
		if err = s.quartz.DeleteJob(key); err != nil {
			return err, false
		}
	}
	//delete(s.nodeStore, id)
	return nil, true
}

func (s *Scheduler) warpJob(jobInfo *model.JobSchedule) (*DbJob, error) {
	if jobInfo == nil {
		return nil, JobScheduleNotFindError
	}

	jobPlug := s.getPlugWithId(jobInfo.Key)
	if jobPlug == nil {
		return nil, JobPlugNotFindError
	}

	// 包装原来逻辑
	warpAction := func(ctx1 context.Context) (res any, err error) {

		jobInfo, _ = s.db.JobSchedule.FindOne(s.ctxDb, jobInfo.Id)
		if jobInfo == nil {
			return nil, JobScheduleNotFindError
		}
		if jobInfo.JobStatus != Enable {
			_, _ = s.deleteJob(jobInfo.Key)
			return nil, nil
		}

		if jobInfo.Parallel == Reject {
			// 非并发, 乐观锁该当前状态
			if !s.db.JobSchedule.LockAndUpdateNowStatusAsRunning(s.ctxDb, jobInfo.Id) {
				return nil, nil
			}
			defer s.db.JobSchedule.UnLockById(s.ctxDb, jobInfo.Id)
		} else {
			if !s.db.JobSchedule.UpdateNowStatusAsRunning(s.ctxDb, jobInfo.Id) {
				return nil, nil
			}
		}

		startTime := time.Now()
		s.setNodeState(jobInfo.Key, Running)

		// 执行原始逻辑
		func() {
			defer func() {
				if r := recover(); r != nil {
					err = errors.New(fmt.Sprint(r))
				}
			}()
			res, err = jobPlug.Action(s.ctx)
		}()

		resStatus := Ok
		if err != nil {
			resStatus = Error
			jobInfo.FailCount = jobInfo.FailCount + 1
		}
		s.setNodeState(jobInfo.Key, Readly)

		jobInfo.NowStatus = Readly
		jobInfo.ExecuteCount = jobInfo.ExecuteCount + 1
		// 当jobInfo.Parallel = "allow"时，如发送并发更新，信息可能不准确
		_ = s.db.JobSchedule.Update(s.ctxDb, jobInfo)

		resultJson := NewResult(res, jobInfo.Key, err, nil).ToJson()
		s.db.InsertRecord(s.ctxDb, jobInfo.Id, startTime, resultJson, string(resStatus), Action)
		return res, err
	}

	job := NewDbJob(warpAction)
	s.setNodeKey(jobInfo.Key, job.Key())
	s.setNodeState(jobInfo.Key, Readly)

	return job, nil
}

func (d *DB) InsertRecord(ctx context.Context, jobId int64, startTime time.Time, result, status string, execType string) (id int64) {
	endTime := time.Now()
	jobRecord := model.JobRecord{
		JobId:     jobId,
		StartTime: &startTime,
		EndTime:   &endTime,
		Result:    result,
		Status:    status,
		UseMilli:  endTime.UnixMilli() - startTime.UnixMilli(),
		ExecType:  execType,
	}
	_ = d.JobRecord.Insert(ctx, &jobRecord)
	return jobRecord.Id
}

func NewResult(Data any, key string, err error, params any) *Result {
	r := &Result{
		Data:   "",
		Key:    key,
		Error:  "",
		Params: "",
	}
	if Data != nil {
		r.Data = Data
	}
	if err != nil {
		r.Error = fmt.Sprint(err)
	}
	if params != nil {
		r.Params = params
	}
	return r
}

func (r *Result) ToJson() string {
	bytes, err := json.Marshal(*r)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func (r *Result) String() string {
	return fmt.Sprint(*r)
}

func VerifyCron(cron string) bool {
	_, err := quartz.NewCronTrigger(cron)
	return err == nil
}
