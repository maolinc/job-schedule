package logic

import (
	"context"
	"job/internal/svc"
	"job/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetScheduleStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetScheduleStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScheduleStatusLogic {
	return &GetScheduleStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 操作Schedule对象，包括restart、stop、start调度
func (l *GetScheduleStatusLogic) GetScheduleStatus(in *pb.ScheduleStatusReq) (*pb.ScheduleStatusResp, error) {
	resp := &pb.ScheduleStatusResp{}
	if l.svcCtx.Scheduler.IsStarted() {
		resp.Status = "running"
	} else {
		resp.Status = "close"
	}

	return resp, nil
}
