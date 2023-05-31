package logic

import (
	"context"
	"github.com/maolinc/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"job/internal/svc"
	"job/internal/tools/errorx"
	"job/internal/tools/job"
	"job/model"
	"job/pb"
)

type UpdateJobScheduleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateJobScheduleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateJobScheduleLogic {
	return &UpdateJobScheduleLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}
func (l *UpdateJobScheduleLogic) UpdateJobSchedule(in *pb.UpdateJobScheduleReq) (resp *pb.UpdateJobScheduleResp, err error) {
	var data model.JobSchedule
	_ = copier.Copiers(&data, in)

	oldJob, _ := l.svcCtx.JobScheduleModel.FindOne(l.ctx, data.Id)
	if oldJob == nil {
		return &pb.UpdateJobScheduleResp{
			ResultStatus: &pb.ResultStatus{
				Code: errorx.RPC_ERROR,
				Msg:  "未找到任务",
			},
		}, nil
	}
	if oldJob.IsLock() {
		return &pb.UpdateJobScheduleResp{
			ResultStatus: &pb.ResultStatus{
				Code: errorx.RPC_ERROR,
				Msg:  "无法修改正在执行的任务",
			},
		}, nil
	}
	if data.Cron != "" && !job.VerifyCron(data.Cron) {
		return &pb.UpdateJobScheduleResp{
			ResultStatus: &pb.ResultStatus{
				Code: errorx.RPC_ERROR,
				Msg:  "cron表达式错误",
			},
		}, nil
	}

	if err, ok := l.svcCtx.JobScheduleModel.UnLockAndUpdate(l.ctx, &data); err != nil || !ok {
		return &pb.UpdateJobScheduleResp{
			ResultStatus: &pb.ResultStatus{
				Code: errorx.RPC_ERROR,
				Msg:  "操作失败",
			},
		}, nil
	}
	// 修改运行中的job
	if oldJob.JobStatus == job.Enable {
		err = l.svcCtx.Scheduler.RestartJob(context.Background(), oldJob.Key)
	} else if data.JobStatus == job.Enable {
		err = l.svcCtx.Scheduler.AddScheduleDbJob(context.Background(), oldJob.Key)
	}

	if err != nil {
		return &pb.UpdateJobScheduleResp{ResultStatus: &pb.ResultStatus{Code: errorx.RPC_ERROR, Msg: "启动失败，请尝试手动重启"}}, nil
	}

	return &pb.UpdateJobScheduleResp{}, nil
}
