package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"job/internal/svc"
	"job/internal/tools/errorx"
	"job/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type JobCompensateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJobCompensateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JobCompensateLogic {
	return &JobCompensateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 补偿任务，立即执行该key对应的任务
func (l *JobCompensateLogic) JobCompensate(in *pb.JobCompensateReq) (*pb.JobCompensateResp, error) {
	params := map[string]any{}
	_ = json.Unmarshal([]byte(in.Params), &params)

	res, err := l.svcCtx.Scheduler.AddCompensate(context.Background(), in.Key, params, in.Background)
	if err != nil {
		return &pb.JobCompensateResp{
			ResultStatus: &pb.ResultStatus{
				Code: errorx.RPC_ERROR,
				Msg:  fmt.Sprint(err),
			},
		}, nil
	}
	result, _ := json.Marshal(res)

	return &pb.JobCompensateResp{Result: string(result)}, nil
}
