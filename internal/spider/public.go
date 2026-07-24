package spider

import "github.com/gogf/gf/v2/frame/g"

const (
	// 待处理任务初始化容量
	PENDING_TASKS_INIT_CAPACITY = 10000
)

var (
	spiderLogger = g.Log("spider")
)

// 任务
type PendingTask struct {
	// 任务Id
	Id string `json:"id"`
	// 是否完成
	IsFinished bool `json:"is_finished"`
}
