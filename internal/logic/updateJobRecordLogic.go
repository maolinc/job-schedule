package logic

import (
	"context"
	"job/internal/svc"
	"job/internal/tools/errorx"
	"job/model"
	"job/pb"

	"github.com/maolinc/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateJobRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateJobRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateJobRecordLogic {
	return &UpdateJobRecordLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}
func (l *UpdateJobRecordLogic) UpdateJobRecord(in *pb.UpdateJobRecordReq) (*pb.UpdateJobRecordResp, error) {
	var data model.JobRecord
	_ = copier.Copiers(&data, in)
	err := l.svcCtx.JobRecordModel.Update(l.ctx, &data)
	if err != nil {
		return &pb.UpdateJobRecordResp{
			ResultStatus: &pb.ResultStatus{
				Code: errorx.RPC_ERROR,
				Msg:  "操作失败",
			},
		}, nil
	}
	return &pb.UpdateJobRecordResp{}, nil
}
