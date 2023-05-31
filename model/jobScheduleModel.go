package model

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"time"
)

const (
	lock   = 1
	unlock = 0
)

var (
	_ JobScheduleModel = (*defaultJobScheduleModel)(nil)
)

type (
	JobSchedule struct {
		Id           int64                 `gorm:"id;primary_key"` //
		JobName      string                `gorm:"job_name"`       //任务名称
		Des          string                `gorm:"des"`            //描述
		Cron         string                `gorm:"cron"`           //调度时间，cron表达式
		ExecuteCount int64                 `gorm:"execute_count"`  //
		FailCount    int64                 `gorm:"fail_count"`     //失败次数
		NowStatus    string                `gorm:"now_status"`     //当前状态，wait | readly | running
		JobStatus    string                `gorm:"job_status"`     //任务转态, stop不加入调度 | enable加入调度中
		DeleteFlag   soft_delete.DeletedAt `gorm:"column:delete_flag;softDelete:flag"`
		Key          string                `gorm:"key"`                               //与程序中的任务key一致，唯一
		Parallel     string                `gorm:"parallel"`                          // 是否允许并发执行，allow | reject
		CreateTime   *time.Time            `gorm:"column:create_time;autoCreateTime"` // 创建时间
		UpdateTime   *time.Time            `gorm:"column:update_time;autoUpdateTime"` // 更新时间
		Lock         int                   `gorm:"column:lock"`                       // 锁，0解锁, 1上锁
	}

	// JobSchedule query cond
	JobScheduleQuery struct {
		SearchBase
		JobSchedule
	}

	JobScheduleModel interface {
		// Trans Transaction
		Trans(ctx context.Context, fn func(ctx context.Context, db *gorm.DB) (err error)) (err error)
		// Builder Custom assembly conditions
		Builder(ctx context.Context, db ...*gorm.DB) (b *gorm.DB)
		Insert(ctx context.Context, data *JobSchedule, db ...*gorm.DB) (err error)
		Update(ctx context.Context, data *JobSchedule, db ...*gorm.DB) (err error)
		// Delete When a tombstone field is set, it is tombstone. Otherwise, it is physically deleted
		Delete(ctx context.Context, id int64, db ...*gorm.DB) (err error)
		// ForceDelete Physical deletion
		ForceDelete(ctx context.Context, id int64, db ...*gorm.DB) (err error)
		Count(ctx context.Context, cond *JobScheduleQuery) (total int64, err error)
		FindOne(ctx context.Context, id int64) (data *JobSchedule, err error)
		// FindByPage Contains total information
		FindByPage(ctx context.Context, cond *JobScheduleQuery) (total int64, list []*JobSchedule, err error)
		// FindListByPage Normal pagination
		FindListByPage(ctx context.Context, cond *JobScheduleQuery) (list []*JobSchedule, err error)
		// FindListByCursor Cursor is required based on cursor paging, Only the primary key is of type int, and other types can be expanded by themselves
		FindListByCursor(ctx context.Context, cond *JobScheduleQuery) (list []*JobSchedule, err error)
		FindAll(ctx context.Context, cond *JobScheduleQuery) (list []*JobSchedule, err error)
		// FindByKey ---------------Write your other interfaces below---------------

		// LockById 上锁
		LockById(ctx context.Context, id int64) bool
		// UnLockById 解锁
		UnLockById(ctx context.Context, id int64) bool
		// LockByKey 上锁
		LockByKey(ctx context.Context, key string) bool
		// UnLockByKey 解锁
		UnLockByKey(ctx context.Context, key string) bool
		// IsLockById 判断锁状态
		IsLockById(id int64) bool
		// IsLockByKey 判断锁状态
		IsLockByKey(key string) bool
		FindByKey(ctx context.Context, key string) (data *JobSchedule, err error)
		// LockAndUpdateNowStatusAsRunning  当未上锁时更新NowStatus为running，并且上锁,  操作失败false  操作成共true
		LockAndUpdateNowStatusAsRunning(ctx context.Context, id int64) bool
		UpdateNowStatusAsRunning(ctx context.Context, id int64) bool
		UpdateNowStatusAsReadly(ctx context.Context, id int64) bool
		// UnLockAndUpdate 当锁处于释放才修改, 修改成功ok
		UnLockAndUpdate(ctx context.Context, data *JobSchedule, db ...*gorm.DB) (err error, ok bool)
	}

	defaultJobScheduleModel struct {
		*customConn
		table string
	}
)

func NewJobScheduleModel(db *gorm.DB) JobScheduleModel {
	return &defaultJobScheduleModel{
		customConn: newCustomConnNoCache(db),
		table:      "job_schedule",
	}
}

func (m *JobSchedule) TableName() string {
	return "`job_schedule`"
}

func (m *defaultJobScheduleModel) conn(ctx context.Context, db ...*gorm.DB) *gorm.DB {
	if len(db) == 0 {
		return m.db.Model(&JobSchedule{}).Session(&gorm.Session{Context: ctx})
	}
	return db[0]
}

func (m *defaultJobScheduleModel) Builder(ctx context.Context, db ...*gorm.DB) (b *gorm.DB) {
	return m.conn(ctx, db...)
}

func (m *defaultJobScheduleModel) Trans(ctx context.Context, fn func(ctx context.Context, db *gorm.DB) (err error)) (err error) {
	return m.conn(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctx, tx)
	})
}

func (m *defaultJobScheduleModel) Insert(ctx context.Context, data *JobSchedule, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Create(data).Error
	})
}

func (m *defaultJobScheduleModel) Update(ctx context.Context, data *JobSchedule, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Model(&data).Updates(data).Error
	})
}

func (m *defaultJobScheduleModel) Delete(ctx context.Context, id int64, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Where(" id = ?", id).Updates(&JobSchedule{DeleteFlag: 1}).Error
	})
}

func (m *defaultJobScheduleModel) ForceDelete(ctx context.Context, id int64, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Unscoped().Where(" id = ?", id).Delete(JobSchedule{}).Error
	})
}

func (m *defaultJobScheduleModel) Count(ctx context.Context, cond *JobScheduleQuery) (total int64, err error) {
	err = m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
	).Where(cond.JobSchedule).Count(&total).Error
	return total, err
}

func (m *defaultJobScheduleModel) FindOne(ctx context.Context, id int64) (data *JobSchedule, err error) {
	err = m.QueryRow(ctx, &data, func(v interface{}) error {
		tx := m.conn(ctx).Where(" id = ?", id).Find(v)
		if tx.RowsAffected == 0 {
			return sql.ErrNoRows
		}
		return tx.Error
	})
	switch err {
	case nil:
		return data, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (m *defaultJobScheduleModel) FindByPage(ctx context.Context, cond *JobScheduleQuery) (total int64, list []*JobSchedule, err error) {
	conn := m.conn(ctx).Scopes(
		orderScope(cond.OrderSort...),
		searchPlusScope(cond.SearchPlus, m.table),
	).Where(cond.JobSchedule)

	total, list, err = pageHandler[*JobSchedule](conn, cond.PageCurrent, cond.PageSize)
	return total, list, err
}

func (m *defaultJobScheduleModel) FindListByPage(ctx context.Context, cond *JobScheduleQuery) (list []*JobSchedule, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
		orderScope(cond.OrderSort...),
		pageScope(cond.PageCurrent, cond.PageSize),
	)
	err = conn.Where(cond.JobSchedule).Find(&list).Error
	return list, err
}

func (m *defaultJobScheduleModel) FindListByCursor(ctx context.Context, cond *JobScheduleQuery) (list []*JobSchedule, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
		cursorScope(cond.Cursor, cond.CursorAsc, cond.PageSize, "id"),
		orderScope(cond.OrderSort...),
	)
	err = conn.Where(cond.JobSchedule).Find(&list).Error
	return list, err
}

func (m *defaultJobScheduleModel) FindAll(ctx context.Context, cond *JobScheduleQuery) (list []*JobSchedule, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
		orderScope(cond.OrderSort...),
	)
	err = conn.Where(cond.JobSchedule).Find(&list).Error
	return list, err
}

func (m *defaultJobScheduleModel) LockById(ctx context.Context, id int64) bool {
	// set lock = 1 where lock = 0 and id = ?   ; >0上锁成功   =0上锁失败
	conn := m.conn(ctx).Scopes(unLockScope())
	affected := conn.Where("`id` = ?", id).Update("`lock`", lock).RowsAffected
	return affected > 0
}

func (m *defaultJobScheduleModel) UnLockById(ctx context.Context, id int64) bool {
	conn := m.conn(ctx)
	affected := conn.Where("`id` = ?", id).Update("`lock`", unlock).RowsAffected
	return affected > 0
}

func (m *defaultJobScheduleModel) LockByKey(ctx context.Context, key string) bool {
	// set lock = 1 where lock = 0 and key = ?   ; >0上锁成功   =0上锁失败
	conn := m.conn(ctx).Scopes(unLockScope())
	affected := conn.Where("`key` = ?", key).Update("`lock`", lock).RowsAffected
	return affected > 0
}

func (m *defaultJobScheduleModel) UnLockByKey(ctx context.Context, key string) bool {
	conn := m.conn(ctx)
	affected := conn.Where("`key` = ?", key).Update("`lock`", unlock).RowsAffected
	return affected > 0
}

func (m *defaultJobScheduleModel) IsLockById(id int64) bool {
	var total int64
	// where lock = 0 and id = ?   ;total>0未上锁   total=0存在锁
	conn := m.conn(context.Background()).Scopes(unLockScope())
	conn.Where("`id` = ?", id).Count(&total)
	return total == 0
}

func (m *defaultJobScheduleModel) IsLockByKey(key string) bool {
	var total int64
	// where lock = 0 and key = ?   ;total>0未上锁   total=0存在锁
	conn := m.conn(context.Background()).Scopes(unLockScope())
	conn.Where("`key` = ?", key).Count(&total)
	return total == 0
}

func (m *defaultJobScheduleModel) FindByKey(ctx context.Context, key string) (data *JobSchedule, err error) {
	err = m.QueryRow(ctx, &data, func(v interface{}) error {
		tx := m.conn(ctx).Where("`key` = ?", key).Find(v)
		if tx.RowsAffected == 0 {
			return sql.ErrNoRows
		}
		return tx.Error
	})
	switch err {
	case nil:
		return data, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (m *defaultJobScheduleModel) LockAndUpdateNowStatusAsRunning(ctx context.Context, id int64) bool {
	conn := m.conn(ctx).Scopes(unLockScope())
	affected := conn.Where("`id` = ?", id).Updates(map[string]interface{}{
		"`lock`":       lock,
		"`now_status`": "running",
	}).RowsAffected
	return affected != 0
}

func (m *defaultJobScheduleModel) UpdateNowStatusAsRunning(ctx context.Context, id int64) bool {
	affected := m.conn(ctx).Where("`id` = ?", id).Update("`now_status`", "running").RowsAffected
	return affected != 0
}

func (m *defaultJobScheduleModel) UpdateNowStatusAsReadly(ctx context.Context, id int64) bool {
	affected := m.conn(ctx).Where("`id` = ?", id).Update("`now_status`", "readly").RowsAffected
	return affected != 0
}

func (m *defaultJobScheduleModel) UnLockAndUpdate(ctx context.Context, data *JobSchedule, db ...*gorm.DB) (err error, ok bool) {
	err = m.Exec(ctx, func() error {
		tx := m.conn(ctx, db...).Scopes(unLockScope()).Model(&data).Updates(data)
		ok = tx.RowsAffected > 0
		return tx.Error
	})
	return err, ok
}

func (m *JobSchedule) IsLock() bool {
	return m.Lock != 0
}

// 判断锁未占用
func unLockScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("`lock` = ?", unlock)
	}
}
