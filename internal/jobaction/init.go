package jobaction

import (
	"context"
	"job/internal/jobaction/jobs"
	"job/internal/svc"
	"job/internal/tools/job"

	"log"
)

// StartJob 加入将定时任务
func StartJob(svc *svc.ServiceContext) []job.Plug {
	plugs := make([]job.Plug, 0)
	ctx := context.Background()

	clsLoginLogCountJob := jobs.NewLoginLogJob(ctx, svc)
	plugs = append(plugs, clsLoginLogCountJob)
	//testDebugAction(clsLoginLogCountJob)

	clsRegisterLogCountJob := jobs.NewSyncOrderJob(ctx, svc)
	plugs = append(plugs, clsRegisterLogCountJob)
	//testDebugAction(clsRegisterLogCountJob)

	errs, ok := svc.Scheduler.Start().Plugs(plugs...)
	if !ok {
		for _, err := range errs {
			log.Println(err)
		}
	}
	return plugs
}

// 启动直接执行action，无需调度，方便debug
func testDebugAction(plug job.Plug) {
	res, err := plug.Action(context.Background())
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
	select {}
}
