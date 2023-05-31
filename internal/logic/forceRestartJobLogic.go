package logic

import (
	"context"
	"job/internal/svc"
	"job/internal/tools/errorx"
	"job/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ForceRestartJobLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewForceRestartJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ForceRestartJobLogic {
	return &ForceRestartJobLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 强制重启job
func (l *ForceRestartJobLogic) ForceRestartJob(in *pb.ForceRestartJobReq) (*pb.ForceRestartJobResp, error) {
	err := l.svcCtx.Scheduler.ForceResetJob(context.Background(), in.Key)
	if err != nil {
		return &pb.ForceRestartJobResp{
			ResultStatus: &pb.ResultStatus{
				Code: errorx.RPC_ERROR,
				Msg:  "重启失败",
			},
		}, nil
	}

	return &pb.ForceRestartJobResp{}, nil
}
