package logic

import (
	"context"
	"job/internal/svc"
	"job/internal/tools/errorx"
	"job/internal/tools/job"
	"job/model"
	"job/pb"

	"github.com/maolinc/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateJobScheduleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateJobScheduleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateJobScheduleLogic {
	return &CreateJobScheduleLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}
func (l *CreateJobScheduleLogic) CreateJobSchedule(in *pb.CreateJobScheduleReq) (*pb.CreateJobScheduleResp, error) {
	data := &model.JobSchedule{}
	_ = copier.Copiers(data, in)
	if data.Key == "" {
		return &pb.CreateJobScheduleResp{
			ResultStatus: &pb.ResultStatus{
				Code: errorx.RPC_ERROR,
				Msg:  "Key必填",
			},
		}, nil
	}

	data.NowStatus = job.Wait
	if data.JobStatus == "" {
		data.JobStatus = job.Stop
	}
	if data.Parallel == "" {
		data.Parallel = job.Reject
	}
	if !job.VerifyCron(data.Cron) {
		return &pb.CreateJobScheduleResp{
			ResultStatus: &pb.ResultStatus{
				Code: errorx.RPC_ERROR,
				Msg:  "Cron表达式错误",
			},
		}, nil
	}

	if j, _ := l.svcCtx.JobScheduleModel.FindByKey(l.ctx, data.Key); j != nil {
		return &pb.CreateJobScheduleResp{
			ResultStatus: &pb.ResultStatus{
				Code: errorx.RPC_ERROR,
				Msg:  "Key已存在",
			},
		}, nil
	}
	err := l.svcCtx.JobScheduleModel.Insert(l.ctx, data)
	if err != nil {
		return &pb.CreateJobScheduleResp{
			ResultStatus: &pb.ResultStatus{
				Code: errorx.RPC_ERROR,
				Msg:  "任务创建失败",
			},
		}, nil
	}

	if data.JobStatus == job.Enable {
		err = l.svcCtx.Scheduler.AddScheduleDbJob(context.Background(), data.Key)
		if err != nil {
			return &pb.CreateJobScheduleResp{
				ResultStatus: &pb.ResultStatus{
					Code: errorx.RPC_ERROR,
					Msg:  "任务创建成功, 启动失败, 请尝试手动启动",
				},
			}, nil
		}
	}

	return &pb.CreateJobScheduleResp{}, nil
}
