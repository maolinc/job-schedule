package logic

import (
	"context"
	"job/internal/svc"
	"job/internal/tools/errorx"
	"job/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	start     = "start"
	restart   = "restart"
	close     = "close"
	forceStop = "forceStop"
)

type OperaScheduleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOperaScheduleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperaScheduleLogic {
	return &OperaScheduleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 操作Schedule对象，包括restart、close、start
func (l *OperaScheduleLogic) OperaSchedule(in *pb.OperaScheduleReq) (*pb.OperaScheduleResp, error) {
	switch in.Action {
	case start:
		l.svcCtx.Scheduler.Start()
	case close:
		l.svcCtx.Scheduler.Stop()
	case restart:
		l.svcCtx.Scheduler.Restart()
	case forceStop:
		l.svcCtx.Scheduler.ForceStop()
	default:
		return &pb.OperaScheduleResp{
			ResultStatus: &pb.ResultStatus{
				Code: errorx.RPC_ERROR,
				Msg:  "未知action",
			},
		}, nil
	}

	resp := &pb.OperaScheduleResp{
		Action: in.Action,
	}
	if l.svcCtx.Scheduler.IsStarted() {
		resp.Status = "running"
	} else {
		resp.Status = "close"
	}

	l.Logger.Infof("OperaSchedule req: %+v, resp: %+v", in, resp)

	return resp, nil
}
