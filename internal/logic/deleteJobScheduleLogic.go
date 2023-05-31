package logic

import (
	"context"
	"job/internal/svc"
	"job/internal/tools/errorx"
	"job/pb"

	"gorm.io/gorm"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteJobScheduleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteJobScheduleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteJobScheduleLogic {
	return &DeleteJobScheduleLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}
func (l *DeleteJobScheduleLogic) DeleteJobSchedule(in *pb.DeleteJobScheduleReq) (*pb.DeleteJobScheduleResp, error) {
	oldJob, _ := l.svcCtx.JobScheduleModel.FindOne(l.ctx, in.Id)
	if oldJob == nil {
		return &pb.DeleteJobScheduleResp{
			ResultStatus: &pb.ResultStatus{
				Code: errorx.RPC_ERROR,
				Msg:  "任务不存在",
			},
		}, nil
	}

	err := l.svcCtx.JobScheduleModel.Trans(l.ctx, func(ctx context.Context, db *gorm.DB) (err error) {
		err = l.svcCtx.JobScheduleModel.Delete(l.ctx, in.Id)
		if err != nil {
			return err
		}
		if err, _ := l.svcCtx.Scheduler.RemoveJob(oldJob.Key); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, errors.Wrapf(errors.New("删除失败"), "DeleteJobSchedule req: %v, error: %v", in, err)
	}

	return &pb.DeleteJobScheduleResp{}, nil
}
