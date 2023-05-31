package job

import "context"

type Plug interface {
	// Action 定时任务需要实现此接口, res的字符数不能大于3000
	Action(ctx context.Context) (res any, err error)
	// Key 定时任务需要实现此接口, 返回的key与数据库的key一致, 每个插件的key必须不同,且不能为""
	Key() string
	// Compensate 补偿接口, res的字符数不能大于3000
	/**
	由于某种原因定任务某次未执行，可主动调用此方法进行补偿
	params执行参数
	*/
	Compensate(ctx context.Context, params map[string]any) (res any, err error)
}
