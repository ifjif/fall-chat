package taskmanager

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"sort"
	"sync"
	"time"
	"zzc/fall-chat/src/utils"
)

// 1.通用任务接口
type Job interface {
	Key() string          // 防抖
	Delay() time.Duration // 延迟时间
	Priority() int        // 数值越小优先级越高（0=紧急 1=普通）
	RetryCount() int      // 当前已重试次数
	MaxRetries() int      // 最大允许重试次数
	Execute() error
	OnComplete()
}

// 2.延迟任务管理器
type DelayedTaskManager struct {
	priorityChans    map[int]chan Job // 不同优先级的独立 Channel
	sortedPriorities []int            // 缓存排序后的优先级，用于reflect.Select
	debounceMap      sync.Map         // 并发安全防抖map
	execChan         chan Job         // 统一的执行通道，保证严格顺序
	wg               sync.WaitGroup   // 优雅关闭时等待
	ctx              context.Context
	cancel           context.CancelFunc
}

var TaskManager *DelayedTaskManager

func init() {
	ctx, cancel := context.WithCancel(context.Background())
	TaskManager = &DelayedTaskManager{
		priorityChans: map[int]chan Job{
			0: make(chan Job, 10000), // 高优先级
			1: make(chan Job, 10000), // 普通
		},
		ctx:      ctx,
		cancel:   cancel,
		execChan: make(chan Job, 10000),
	}

	// 初始化排序后的优先级列表
	for p := range TaskManager.priorityChans {
		TaskManager.sortedPriorities = append(TaskManager.sortedPriorities, p)
	}
	sort.Ints(TaskManager.sortedPriorities)

	// 启动执行器
	go TaskManager.StartExecuter()
	// 启动调度
	go TaskManager.StartDispatcher()
}

// 3. 执行器，（负责 接收任务、防抖、触发执行)
func (m *DelayedTaskManager) StartExecuter() {
	for {
		select {
		case job, ok := <-m.execChan:
			if !ok {
				utils.LogInfo("TaskManager execute channel closed")
				return
			}
			m.handleJob(job)
		case <-m.ctx.Done():
			return
		}
	}
}

// 3.调度器（使用reflect select 动态监听多通道）
func (m *DelayedTaskManager) StartDispatcher() {
	defer m.wg.Wait() // 优雅关闭时，确保dispatcher退出
	// 动态构建 reflect.Select 的cases
	// 这里将高优先级channel放在前面
	cases := make([]reflect.SelectCase, len(m.sortedPriorities))
	for i, p := range m.sortedPriorities {
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(m.priorityChans[p]),
		}
	}

	// 添加最后一个case 用于监听 Context 取消信号
	allCases := append(cases, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(m.ctx.Done()),
	})

	for {
		var job Job
		var ok bool
		var found bool
		// 1.尝试非阻塞地按优先级顺序提取任务（0,1,2,3...）
		for i := range cases {
			chosen, value, recvOk := reflect.Select([]reflect.SelectCase{
				cases[i],
				{Dir: reflect.SelectDefault},
			})

			if chosen == 0 && recvOk {
				job = value.Interface().(Job)
				ok = true
				found = true
				break
			}
		}

		// 2.所有优先级通道都没有数据，则阻塞等待任意通道
		if !found {
			// 阻塞等待，直到有任务进来或收到关闭信号
			chosen, value, recvOk := reflect.Select(allCases)

			// 处理Context 取消(优雅关闭)
			if chosen == len(allCases)-1 {
				utils.LogInfo("TaskManager Dispatcher context canceled.")
				m.cleanupDebounceMap()
				return
			}

			// 处理channel 关闭
			if !recvOk {
				utils.LogInfo("Priority channel %d closed", m.sortedPriorities[chosen])
				return
			}

			// 正常接受到任务
			job = value.Interface().(Job)
			ok = true
		}

		// 3.将拿到的 放入统一的执行通道
		if ok {
			select {
			case m.execChan <- job:
			case <-m.ctx.Done():
				return
			}
		}
	}
}

// 内部方法：处理具体的防抖逻辑
func (m *DelayedTaskManager) handleJob(job Job) {
	key := job.Key()

	// 防抖逻辑，取消旧定时任务，设置新定时任务
	if val, ok := m.debounceMap.Load(key); ok {
		if timer, ok := val.(*time.Timer); ok {
			timer.Stop()
		}
	}

	// 设置新的延迟任务
	timer := time.AfterFunc(job.Delay(), func() {
		go m.safeExecute(job)
		m.debounceMap.Delete(key)
	})

	m.debounceMap.Store(key, timer)
}

// 安全执行器包装：捕获panic 并处理重试
func (m *DelayedTaskManager) safeExecute(job Job) {
	var err error

	// 捕获panic
	defer func() {
		if r := recover(); r != nil {
			utils.LogWarn("Job %s panicked: %v", job.Key(), r)
			m.handleRetry(job, fmt.Errorf("panic: %v", r))
		}
	}()

	err = job.Execute()
	if err != nil {
		m.handleRetry(job, err)
	} else {
		job.OnComplete()
	}
}

// 重试逻辑（指数退避）
func (m *DelayedTaskManager) handleRetry(job Job, err error) {
	currentRetry := job.RetryCount()
	maxRetry := job.MaxRetries()

	// 执行完后 +1, 如果初始是0, execute后就是1
	if currentRetry > maxRetry {
		utils.LogWarn("Job %s failed after %d retries. Last error: %v. Dropping job", job.Key(), maxRetry, err)
		return
	}

	// 指数退避时间：基础 1s , 最大不超过 30s
	backOff := time.Duration(math.Pow(2, float64(currentRetry))) * time.Second
	if backOff > 30*time.Second {
		backOff = 30 * time.Second
	}

	utils.LogWarn("Job %s failed (attempt %d/%d). Retrying in %v. Error: %v",
		job.Key(), currentRetry, maxRetry, backOff, err)

	// 重新提交到channel中
	go func() {
		time.Sleep(backOff)
		m.Submit(job)
	}()
}

// 4.安全的任务提交方法
func (m *DelayedTaskManager) Submit(job Job) {
	p := job.Priority()
	ch, exists := m.priorityChans[p]

	if !exists {
		utils.LogWarn("Unknown priority %d for job: %s", p, job.Key())
		return
	}

	select {
	case ch <- job:
	default:
		utils.LogWarn("TaskManager channel is full, dropping job: %s", job.Key())
	}
}

// 5.优雅关闭
func (m *DelayedTaskManager) Shutdown() {
	utils.LogWarn("shutting down TaskManager...")
	m.cancel() // 发送停止信号
	// 关闭所有 优先级 channel
	for _, ch := range m.priorityChans {
		close(ch)
	}
	// 关闭执行通道
	close(m.execChan)
	// 不需要 wg.wait, 因为是单协程
}

// 清理防抖map中的 定时器
func (m *DelayedTaskManager) cleanupDebounceMap() {
	m.debounceMap.Range(func(key, value any) bool {
		if timer, ok := value.(*time.Timer); ok {
			timer.Stop()
		}
		m.debounceMap.Delete(key)
		return true
	})
}
