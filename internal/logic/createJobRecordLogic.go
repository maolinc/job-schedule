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

type CreateJobRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateJobRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateJobRecordLogic {
	return &CreateJobRecordLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}
func (l *CreateJobRecordLogic) CreateJobRecord(in *pb.CreateJobRecordReq) (*pb.CreateJobRecordResp, error) {
	data := &model.JobRecord{}
	_ = copier.Copiers(data, in)
	err := l.svcCtx.JobRecordModel.Insert(l.ctx, data)
	if err != nil {
		return &pb.CreateJobRecordResp{
			ResultStatus: &pb.ResultStatus{
				Code: errorx.RPC_ERROR,
				Msg:  "操作失败",
			},
		}, nil
	}
	return &pb.CreateJobRecordResp{}, nil
}
