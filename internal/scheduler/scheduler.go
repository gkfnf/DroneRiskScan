package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// TaskScheduler 任务调度器接口
type TaskScheduler interface {
	Start(ctx context.Context) error
	Stop() error
	Submit(task *Task, handler TaskHandler) error
	GetStatus() *SchedulerStatus
}

// TaskHandler 任务处理函数
type TaskHandler func(ctx context.Context, task *Task) error

// Priority 任务优先级
type Priority int

const (
	PriorityLow Priority = iota
	PriorityNormal
	PriorityHigh
	PriorityCritical
)

// TaskType 任务类型
type TaskType string

const (
	TaskTypeScan     TaskType = "scan"
	TaskTypeReport   TaskType = "report"
	TaskTypeValidate TaskType = "validate"
)

// TaskStatus 任务状态
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
	TaskStatusCancelled TaskStatus = "cancelled"
)

// Task 任务定义
type Task struct {
	ID        string                 `json:"id"`
	Type      TaskType               `json:"type"`
	Priority  Priority               `json:"priority"`
	Status    TaskStatus             `json:"status"`
	Payload   map[string]interface{} `json:"payload"`
	CreatedAt time.Time              `json:"created_at"`
	StartedAt *time.Time             `json:"started_at,omitempty"`
	EndedAt   *time.Time             `json:"ended_at,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Retries   int                    `json:"retries"`
	MaxRetries int                   `json:"max_retries"`
}

// Config 调度器配置
type Config struct {
	MaxWorkers    int
	QueueSize     int
	RetryAttempts int
	RetryDelay    time.Duration
}

// SchedulerStatus 调度器状态
type SchedulerStatus struct {
	IsRunning     bool              `json:"is_running"`
	ActiveWorkers int               `json:"active_workers"`
	QueuedTasks   int               `json:"queued_tasks"`
	CompletedTasks int              `json:"completed_tasks"`
	FailedTasks   int               `json:"failed_tasks"`
	TaskStats     map[TaskType]int  `json:"task_stats"`
}

// DefaultTaskScheduler 默认任务调度器实现
type DefaultTaskScheduler struct {
	config       *Config
	workers      []*worker
	taskQueue    chan *taskItem
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	mutex        sync.RWMutex
	status       *SchedulerStatus
	isRunning    bool
}

// taskItem 任务项（包含任务和处理函数）
type taskItem struct {
	task    *Task
	handler TaskHandler
}

// worker 工作者
type worker struct {
	id       int
	taskChan chan *taskItem
	quit     chan bool
	wg       *sync.WaitGroup
}

// NewTaskScheduler 创建新的任务调度器
func NewTaskScheduler(config *Config) TaskScheduler {
	if config == nil {
		config = &Config{
			MaxWorkers:    5,
			QueueSize:     100,
			RetryAttempts: 3,
			RetryDelay:    time.Second,
		}
	}

	return &DefaultTaskScheduler{
		config:    config,
		taskQueue: make(chan *taskItem, config.QueueSize),
		status: &SchedulerStatus{
			TaskStats: make(map[TaskType]int),
		},
	}
}

// Start 启动调度器
func (ts *DefaultTaskScheduler) Start(ctx context.Context) error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	if ts.isRunning {
		return fmt.Errorf("调度器已在运行中")
	}

	ts.ctx, ts.cancel = context.WithCancel(ctx)
	ts.workers = make([]*worker, ts.config.MaxWorkers)
	ts.status.IsRunning = true
	ts.isRunning = true

	// 启动工作者
	for i := 0; i < ts.config.MaxWorkers; i++ {
		worker := &worker{
			id:       i,
			taskChan: make(chan *taskItem),
			quit:     make(chan bool),
			wg:       &ts.wg,
		}
		ts.workers[i] = worker
		
		ts.wg.Add(1)
		go ts.runWorker(worker)
	}

	// 启动任务分发器
	ts.wg.Add(1)
	go ts.runDispatcher()

	return nil
}

// Stop 停止调度器
func (ts *DefaultTaskScheduler) Stop() error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	if !ts.isRunning {
		return nil
	}

	// 取消上下文
	if ts.cancel != nil {
		ts.cancel()
	}

	// 停止所有工作者
	for _, worker := range ts.workers {
		close(worker.quit)
	}

	// 关闭任务队列
	close(ts.taskQueue)

	// 等待所有工作者完成
	ts.wg.Wait()

	ts.status.IsRunning = false
	ts.status.ActiveWorkers = 0
	ts.isRunning = false

	return nil
}

// Submit 提交任务
func (ts *DefaultTaskScheduler) Submit(task *Task, handler TaskHandler) error {
	if !ts.isRunning {
		return fmt.Errorf("调度器未运行")
	}

	if task == nil {
		return fmt.Errorf("任务不能为空")
	}

	if handler == nil {
		return fmt.Errorf("任务处理函数不能为空")
	}

	// 设置任务默认值
	if task.Status == "" {
		task.Status = TaskStatusPending
	}
	if task.MaxRetries == 0 {
		task.MaxRetries = ts.config.RetryAttempts
	}
	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now()
	}

	// 创建任务项
	item := &taskItem{
		task:    task,
		handler: handler,
	}

	// 提交到队列
	select {
	case ts.taskQueue <- item:
		ts.mutex.Lock()
		ts.status.QueuedTasks++
		ts.status.TaskStats[task.Type]++
		ts.mutex.Unlock()
		return nil
	case <-ts.ctx.Done():
		return fmt.Errorf("调度器已停止")
	default:
		return fmt.Errorf("任务队列已满")
	}
}

// GetStatus 获取调度器状态
func (ts *DefaultTaskScheduler) GetStatus() *SchedulerStatus {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()

	// 创建状态副本
	status := &SchedulerStatus{
		IsRunning:      ts.status.IsRunning,
		ActiveWorkers:  ts.status.ActiveWorkers,
		QueuedTasks:    len(ts.taskQueue),
		CompletedTasks: ts.status.CompletedTasks,
		FailedTasks:    ts.status.FailedTasks,
		TaskStats:      make(map[TaskType]int),
	}

	// 复制任务统计
	for taskType, count := range ts.status.TaskStats {
		status.TaskStats[taskType] = count
	}

	return status
}

// runDispatcher 运行任务分发器
func (ts *DefaultTaskScheduler) runDispatcher() {
	defer ts.wg.Done()

	workerIndex := 0
	for {
		select {
		case item, ok := <-ts.taskQueue:
			if !ok {
				return // 队列已关闭
			}

			// 选择工作者（简单的轮询分配）
			worker := ts.workers[workerIndex]
			workerIndex = (workerIndex + 1) % len(ts.workers)

			// 发送任务给工作者
			select {
			case worker.taskChan <- item:
				ts.mutex.Lock()
				ts.status.QueuedTasks--
				ts.mutex.Unlock()
			case <-ts.ctx.Done():
				return
			}

		case <-ts.ctx.Done():
			return
		}
	}
}

// runWorker 运行工作者
func (ts *DefaultTaskScheduler) runWorker(w *worker) {
	defer ts.wg.Done()

	for {
		select {
		case item, ok := <-w.taskChan:
			if !ok {
				return
			}

			ts.mutex.Lock()
			ts.status.ActiveWorkers++
			ts.mutex.Unlock()

			// 执行任务
			ts.executeTask(w, item)

			ts.mutex.Lock()
			ts.status.ActiveWorkers--
			ts.mutex.Unlock()

		case <-w.quit:
			return
		case <-ts.ctx.Done():
			return
		}
	}
}

// executeTask 执行任务
func (ts *DefaultTaskScheduler) executeTask(w *worker, item *taskItem) {
	task := item.task
	handler := item.handler

	// 更新任务状态
	task.Status = TaskStatusRunning
	now := time.Now()
	task.StartedAt = &now

	// 创建任务上下文
	taskCtx, cancel := context.WithCancel(ts.ctx)
	defer cancel()

	// 执行任务处理函数
	err := handler(taskCtx, task)

	// 更新任务结束时间
	endTime := time.Now()
	task.EndedAt = &endTime

	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	if err != nil {
		task.Error = err.Error()
		
		// 检查是否需要重试
		if task.Retries < task.MaxRetries {
			task.Retries++
			task.Status = TaskStatusPending
			
			// 延迟后重新提交任务
			go func() {
				time.Sleep(ts.config.RetryDelay)
				select {
				case ts.taskQueue <- item:
					ts.status.QueuedTasks++
				case <-ts.ctx.Done():
				}
			}()
			return
		}

		task.Status = TaskStatusFailed
		ts.status.FailedTasks++
	} else {
		task.Status = TaskStatusCompleted
		ts.status.CompletedTasks++
	}
}

// PriorityTaskScheduler 支持优先级的任务调度器
type PriorityTaskScheduler struct {
	*DefaultTaskScheduler
	priorityQueues map[Priority]chan *taskItem
}

// NewPriorityTaskScheduler 创建支持优先级的任务调度器
func NewPriorityTaskScheduler(config *Config) TaskScheduler {
	base := NewTaskScheduler(config).(*DefaultTaskScheduler)
	
	pts := &PriorityTaskScheduler{
		DefaultTaskScheduler: base,
		priorityQueues: map[Priority]chan *taskItem{
			PriorityCritical: make(chan *taskItem, config.QueueSize/4),
			PriorityHigh:     make(chan *taskItem, config.QueueSize/4),
			PriorityNormal:   make(chan *taskItem, config.QueueSize/2),
			PriorityLow:      make(chan *taskItem, config.QueueSize/4),
		},
	}

	return pts
}

// Submit 提交任务到对应优先级队列
func (pts *PriorityTaskScheduler) Submit(task *Task, handler TaskHandler) error {
	if !pts.isRunning {
		return fmt.Errorf("调度器未运行")
	}

	queue, exists := pts.priorityQueues[task.Priority]
	if !exists {
		task.Priority = PriorityNormal
		queue = pts.priorityQueues[PriorityNormal]
	}

	item := &taskItem{
		task:    task,
		handler: handler,
	}

	select {
	case queue <- item:
		pts.mutex.Lock()
		pts.status.QueuedTasks++
		pts.status.TaskStats[task.Type]++
		pts.mutex.Unlock()
		return nil
	case <-pts.ctx.Done():
		return fmt.Errorf("调度器已停止")
	default:
		return fmt.Errorf("优先级队列已满")
	}
}