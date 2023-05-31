package jobs

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"job/internal/svc"
	"job/internal/tools/datex"
	"job/internal/tools/job"

	"time"
)

const (
	syncOrderJobKey = "SyncOrderJob"
)

var (
	_ job.Plug = (*SyncOrderJob)(nil)
)

type SyncOrderJob struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncOrderJob(ctx context.Context, svcCtx *svc.ServiceContext) *SyncOrderJob {
	return &SyncOrderJob{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (c *SyncOrderJob) Action(ctx context.Context) (res any, err error) {
	end := time.Now()
	start := time.Now().AddDate(0, 0, -1)
	return c.ComputerRegisterLog(ctx, start, end)
}

func (c *SyncOrderJob) Key() string {
	return syncOrderJobKey
}

func (c *SyncOrderJob) Compensate(ctx context.Context, params map[string]interface{}) (res any, err error) {
	var (
		start time.Time
		end   time.Time
	)

	start, end, err = parseParamsTime(params)
	if err != nil {
		return nil, err
	}

	return c.ComputerRegisterLog(ctx, start, end)
}

func (c *SyncOrderJob) ComputerRegisterLog(ctx context.Context, start, end time.Time) (res any, err error) {
	// todo logic
	res = fmt.Sprintf("同步%s至%s的订单：共处理%d条订单", datex.FormatDateTime(start), datex.FormatDateTime(end), 20000)

	return res, err
}
