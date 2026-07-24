package spider

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sync"
	"time"
)

const (
	// 进度文件名
	PROGRESS_FILE_NAME = "progress.json"
)

var (
	progressManager     *ProgressManager
	progressManagerOnce sync.Once
)

// 进度管理器
type ProgressManager struct {
	// 锁
	mu sync.RWMutex
	// 保存进度的文件路径
	filePath string
	// 是否完成
	IsDone bool `json:"is_done"`
	// 完成时间戳
	Timestamp int64 `json:"timestamp"`
	// 待处理的任务
	Tasks []PendingTask `json:"tasks"`
}

func NewProgressManager(dir string) *ProgressManager {
	progressManagerOnce.Do(func() {
		progressManager = &ProgressManager{
			filePath: path.Join(dir, PROGRESS_FILE_NAME),
			Tasks:    make([]PendingTask, 0, PENDING_TASKS_INIT_CAPACITY),
		}
	})
	return progressManager
}

// 添加任务
func (pm *ProgressManager) PutTask(task PendingTask) bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	exist := false
	for _, t := range pm.Tasks {
		if t.Id == task.Id {
			exist = true
			break
		}
	}
	// 不存在则放入
	if !exist {
		pm.Tasks = append(pm.Tasks, task)
	}

	return exist
}

// 获取未执行的任务
func (pm *ProgressManager) UnexecutedTasks() (tasks []PendingTask) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	for _, task := range pm.Tasks {
		if !task.IsFinished {
			tasks = append(tasks, task)
		}
	}
	return
}

// 清空所有任务
func (pm *ProgressManager) ClearTasks() {
	pm.Tasks = make([]PendingTask, 0, PENDING_TASKS_INIT_CAPACITY)
	pm.IsDone = false
	pm.Timestamp = 0
}

// 获取任务状态
func (pm *ProgressManager) TaskStatus(ctx context.Context, id string) (bool, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	for _, task := range pm.Tasks {
		if task.Id == id {
			return task.IsFinished, nil
		}
	}

	return false, fmt.Errorf("task %s is not on the tasks list", id)
}

// 标记任务状态
func (pm *ProgressManager) MarkTask(ctx context.Context, id string, isFinished bool) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	finishedCnt := 0
	idx := -1
	for i, task := range pm.Tasks {
		if task.Id == id {
			idx = i
		}
		if task.IsFinished {
			finishedCnt++
		}
	}
	if idx == -1 {
		spiderLogger.Warningf(ctx, "task %s is not on the tasks list", id)
	} else {
		pm.Tasks[idx].IsFinished = isFinished
		finishedCnt++
	}
	if finishedCnt == len(pm.Tasks) {
		pm.IsDone = true
		pm.Timestamp = time.Now().Unix()
	}
}

// 是否处理完所有任务
func (pm *ProgressManager) Done() bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	return pm.IsDone
}

// 处理完成
func (pm *ProgressManager) SetDone() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.IsDone = true
}

// 上次完成时间戳
func (pm *ProgressManager) LastTS() int64 {
	return pm.Timestamp
}

// 加载进度
func (pm *ProgressManager) Load(ctx context.Context) (err error) {
	// 存在进度文件则加载
	if _, err = os.Stat(pm.filePath); err == nil {
		spiderLogger.Debugf(ctx, "found process file %s and load", PROGRESS_FILE_NAME)
		byteValues, rErr := os.ReadFile(pm.filePath)
		if rErr != nil {
			return rErr
		}
		err = json.Unmarshal(byteValues, pm)
		if err == nil {
			spiderLogger.Debugf(ctx, "process file %s load success", PROGRESS_FILE_NAME)
		}
		return
	} else {
		spiderLogger.Debugf(ctx, "process file %s not found skip load", PROGRESS_FILE_NAME)
		return nil
	}
}

// 保存进度
func (pm *ProgressManager) Save(ctx context.Context) (err error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	byteValues, err := json.Marshal(pm)
	if err != nil {
		return
	}
	err = os.WriteFile(pm.filePath, byteValues, 0666)
	return
}
