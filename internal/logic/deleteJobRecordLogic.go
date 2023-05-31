package logic

import (
	"context"
	"job/internal/svc"
	"job/internal/tools/errorx"
	"job/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteJobRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteJobRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteJobRecordLogic {
	return &DeleteJobRecordLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}
func (l *DeleteJobRecordLogic) DeleteJobRecord(in *pb.DeleteJobRecordReq) (*pb.DeleteJobRecordResp, error) {
	err := l.svcCtx.JobRecordModel.Delete(l.ctx, in.Id)
	if err != nil {
		return &pb.DeleteJobRecordResp{
			ResultStatus: &pb.ResultStatus{
				Code: errorx.RPC_ERROR,
				Msg:  "操作失败",
			},
		}, nil
	}
	return &pb.DeleteJobRecordResp{}, nil
}
