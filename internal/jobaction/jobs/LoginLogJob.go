package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"job/internal/svc"
	"job/internal/tools/datex"
	"job/internal/tools/job"

	"time"
)

const (
	loginLogJobKey = "LoginLogJob"
)

var (
	_         job.Plug = (*LoginLogJob)(nil) // 必须实现plug的方法，否则panic
	paramsErr          = errors.New("缺少参数,参数格式: {startTime: string, endTime: string} \n参数描述: \nstartTime: 开始时间,格式YYYY-MM-DD hh:mm:ss; \nendTime: 结束时间,格式YYYY-MM-DD hh:mm:ss")
)

type LoginLogJob struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogJob(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogJob {
	return &LoginLogJob{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (c *LoginLogJob) Action(ctx context.Context) (res any, err error) {
	end := time.Now()
	start := time.Now().AddDate(0, 0, -1)
	return c.ComputerLoginLog(ctx, start, end)
}

func (c *LoginLogJob) Key() string {
	return loginLogJobKey
}

// Compensate 当定时任务某次未执行，可执行此方法作为补偿
func (c *LoginLogJob) Compensate(ctx context.Context, params map[string]any) (res any, err error) {
	var (
		start time.Time
		end   time.Time
	)

	// 补偿方法首先解析参数
	start, end, err = parseParamsTime(params)
	if err != nil {
		return nil, err
	}

	return c.ComputerLoginLog(ctx, start, end)
}

func parseParamsTime(params map[string]any) (start, end time.Time, err error) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = paramsErr
			}
		}()
		startUnix := params["startTime"].(string)
		endUnix := params["endTime"].(string)
		if start, err = datex.ParseDateTimeLocation(startUnix); err != nil {
			return
		}
		if end, err = datex.ParseDateTimeLocation(endUnix); err != nil {
			return
		}
	}()
	if err != nil {
		return time.Now(), time.Now(), err
	}
	return start, end, nil
}

func (c *LoginLogJob) ComputerLoginLog(ctx context.Context, start, end time.Time) (res any, err error) {
	// todo logic
	res = fmt.Sprintf("处理%s至%s登录日志：循环%d次，共处理%d条, 留存计算完成", datex.FormatDateTime(start), datex.FormatDateTime(end), 100, 1000000)

	return res, err
}

func getLastTs(lastLog string) string {
	source := map[string]interface{}{}
	err := json.Unmarshal([]byte(lastLog), &source)
	if err != nil {
		return "0"
	}
	if ts, ok := source["ts"]; ok {
		return ts.(string)
	}
	return "0"
}
