package logic

import (
	"context"
	"job/internal/svc"
	"job/internal/tools/errorx"
	"job/pb"

	"github.com/maolinc/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type DetailJobRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDetailJobRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailJobRecordLogic {
	return &DetailJobRecordLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}
func (l *DetailJobRecordLogic) DetailJobRecord(in *pb.DetailJobRecordReq) (*pb.DetailJobRecordResp, error) {
	data, err := l.svcCtx.JobRecordModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return &pb.DetailJobRecordResp{
			ResultStatus: &pb.ResultStatus{
				Code: errorx.RPC_ERROR,
				Msg:  "查询失败",
			},
		}, nil
	}
	if data == nil {
		return &pb.DetailJobRecordResp{
			ResultStatus: &pb.ResultStatus{
				Code: errorx.RPC_ERROR,
				Msg:  "资源不存在",
			},
		}, nil
	}
	resp := &pb.DetailJobRecordResp{}
	_ = copier.Copiers(&resp, data)
	return resp, nil
}
