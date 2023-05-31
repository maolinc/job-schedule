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

type PageJobScheduleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPageJobScheduleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageJobScheduleLogic {
	return &PageJobScheduleLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}
func (l *PageJobScheduleLogic) PageJobSchedule(in *pb.SearchJobScheduleReq) (*pb.SearchJobScheduleResp, error) {
	var cond model.JobScheduleQuery
	_ = copier.Copiers(&cond, in)
	_ = json.Unmarshal([]byte(in.SearchPlus), &cond.SearchPlus)

	total, list, err := l.svcCtx.JobScheduleModel.FindByPage(l.ctx, &cond)
	if err != nil {
		return &pb.SearchJobScheduleResp{
			ResultStatus: &pb.ResultStatus{
				Code: errorx.RPC_ERROR,
				Msg:  "查询失败",
			},
		}, nil
	}
	resp := &pb.SearchJobScheduleResp{}
	resp.Total = total
	_ = copier.Copiers(&resp.List, list)
	return resp, nil
}
