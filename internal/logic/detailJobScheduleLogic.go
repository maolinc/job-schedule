package logic

import (
	"context"
	"job/internal/svc"
	"job/internal/tools/errorx"
	"job/pb"

	"github.com/maolinc/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type DetailJobScheduleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDetailJobScheduleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailJobScheduleLogic {
	return &DetailJobScheduleLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}
func (l *DetailJobScheduleLogic) DetailJobSchedule(in *pb.DetailJobScheduleReq) (*pb.DetailJobScheduleResp, error) {
	data, err := l.svcCtx.JobScheduleModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return &pb.DetailJobScheduleResp{
			ResultStatus: &pb.ResultStatus{
				Code: errorx.RPC_ERROR,
				Msg:  "查询失败",
			},
		}, nil
	}
	if data == nil {
		return &pb.DetailJobScheduleResp{
			ResultStatus: &pb.ResultStatus{
				Code: errorx.RPC_ERROR,
				Msg:  "资源不存在",
			},
		}, nil
	}
	resp := &pb.DetailJobScheduleResp{}
	_ = copier.Copiers(&resp, data)
	return resp, nil
}
