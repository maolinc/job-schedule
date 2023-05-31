package logic

import (
	"context"
	"encoding/json"
	"job/internal/svc"
	"job/internal/tools/errorx"
	"job/model"
	"job/pb"

	"github.com/maolinc/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type PageJobRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPageJobRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageJobRecordLogic {
	return &PageJobRecordLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}
func (l *PageJobRecordLogic) PageJobRecord(in *pb.SearchJobRecordReq) (*pb.SearchJobRecordResp, error) {
	var cond model.JobRecordQuery
	_ = copier.Copiers(&cond, in)
	_ = json.Unmarshal([]byte(in.SearchPlus), &cond.SearchPlus)

	total, list, err := l.svcCtx.JobRecordModel.FindByPage(l.ctx, &cond)
	if err != nil {
		return &pb.SearchJobRecordResp{
			ResultStatus: &pb.ResultStatus{
				Code: errorx.RPC_ERROR,
				Msg:  "查询失败",
			},
		}, nil
	}
	resp := &pb.SearchJobRecordResp{}
	resp.Total = total
	_ = copier.Copiers(&resp.List, list)
	return resp, nil
}
