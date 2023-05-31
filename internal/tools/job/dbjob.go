package job

import (
	"context"
	"fmt"
	"github.com/reugn/go-quartz/quartz"
)

type DbJobStatus string

const (
	Na    DbJobStatus = "na"
	Ok    DbJobStatus = "ok"
	Error DbJobStatus = "error"
)

var (
	_ quartz.Job = (*DbJob)(nil)
)

type Function func(context.Context) (res any, err error)

type DbJob struct {
	function  *Function
	desc      string
	Result    any
	Error     error
	JobStatus DbJobStatus
}

func NewDbJob(function Function) *DbJob {
	return &DbJob{
		function:  &function,
		desc:      fmt.Sprintf("DbJob:%p", &function),
		Result:    nil,
		Error:     nil,
		JobStatus: Na,
	}
}

func (j *DbJob) Execute(ctx context.Context) {
	// 兼容在任务中移除任务后  j为nil
	if j == nil {
		return
	}
	result, err := (*j.function)(ctx)
	if err != nil {
		j.JobStatus = Error
		j.Result = nil
		j.Error = err
	} else {
		j.JobStatus = Ok
		j.Error = nil
		j.Result = &result
	}
}

func (j *DbJob) Description() string {
	return j.desc
}

func (j *DbJob) Key() int {
	return quartz.HashCode(fmt.Sprintf("%s:%p", j.desc, j.function))
}
