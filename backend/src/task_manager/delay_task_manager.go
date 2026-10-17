package taskmanager

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"sync"
	"time"
	"zzc/fall-chat/src/utils"
)

// 防抖任务管理器
// - 延迟执行
// - 后加入的任务可能取代前面的任务
// - 任务加入不同的优先级通道，然后统一交给一个执行通道

// 1.通用任务接口
type Job interface {
	Key() string                // 任务唯一键
	Delay() time.Duration       // 延迟执行时间
	Priority() int              // 任务的优先级,越小越高(0=紧急 1=普通)
	Execute() error             // 任务逻辑
	ShouldReplace(old Job) bool // 是否要取代前面的任务, false则被丢弃
	OnComplete()                // 任务完成后执行
	RetryCount() int            // 重试次数
	MaxRetries() int            // 最大重试次数
}

// 2. 防抖实体
type debounceEntry struct {
	job   Job
	timer *time.Timer
}

// 3. 任务管理器
type DelayedTaskManager struct {
	priorityChans    map[int]chan Job // 优先级chan
	sortedPriorities []int            // 缓存排序后的优先级，用于reflect.select
	execChan         chan Job         // 执行通道
	debounceMap      sync.Map         // 并发安全防抖map
	wg               sync.WaitGroup   // 优雅关闭时等待
	ctx              context.Context
	cancel           context.CancelFunc
}

// 执行器（接收任务，防抖，执行）
func (m *DelayedTaskManager) StartExecuter() {
	for {
		select {
		case job, ok := <-m.execChan:
			if !ok {
				utils.LogWarn("TaskManager execute channel closed")
				return
			}
			m.handleJob(job)
		case <-m.ctx.Done():
			return
		}
	}
}

// 由执行器调用
func (m *DelayedTaskManager) handleJob(job Job) {
	key := job.Key()

	// 防抖，取消旧任务，执行新任务
	if val, ok := m.debounceMap.Load(key); ok {
		if entry, ok := val.(*debounceEntry); ok {
			// 判断是否允许覆盖
			if !job.ShouldReplace(entry.job) {
				utils.LogWarn("Job %s rejected by ShouldReplace condition, dropping.", key)
				return
			}
			// 停止旧定时器
			entry.timer.Stop()
		}
	}

	// 设置新定时任务
	timer := time.AfterFunc(job.Delay(), func() {
		go m.safeExecute(job)
		m.debounceMap.Delete(key)
	})

	m.debounceMap.Store(key, &debounceEntry{
		job:   job,
		timer: timer,
	})
}

// 安全执行任务，捕获panic, 执行重试
func (m *DelayedTaskManager) safeExecute(job Job) {
	var err error

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

// 重试逻辑，指数退避
func (m *DelayedTaskManager) handleRetry(job Job, err error) {
	currentRetry := job.RetryCount()
	maxRetry := job.MaxRetries()

	// execute后+1
	if currentRetry > maxRetry {
		utils.LogWarn("job %s failed after %d retries. Last error %v. Dropping job.", job.Key(), maxRetry, err)
		return
	}

	// 指数退避时间,基础 1s 最大30s
	backOff := time.Duration(math.Pow(2, float64(currentRetry))) * time.Second
	if backOff > 30*time.Second {
		backOff = 30 * time.Second
	}

	utils.LogWarn("job %s failed. attempt(%d/%d). Retrying in %v. Error: %v",
		job.Key(), currentRetry, maxRetry, backOff, err)

	// 重新提交到 任务管理器
	go func() {
		time.Sleep(backOff)
		m.Submit(job)
	}()
}

// 调度器（使用reflect.Select 动态监听多通道）
func (m *DelayedTaskManager) StartDispatcher() {
	defer m.wg.Wait()

	// 动态构建Select, 高优先级在最前面
	cases := make([]reflect.SelectCase, len(m.sortedPriorities))
	for i, p := range m.sortedPriorities {
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(m.priorityChans[p]),
		}
	}

	// 添加最后一个 case 来监听context
	allCases := append(cases, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(m.ctx.Done()),
	})

	// 先对所有优先级单独 非阻塞的select
	// 如果都没有找到，整体阻塞select
	for {
		var job Job
		var ok bool    // 最终是否存在
		var found bool // 是否从非阻塞中找到

		for i := range cases {
			chosen, value, recvOk := reflect.Select([]reflect.SelectCase{
				cases[i],
				{
					Dir: reflect.SelectDefault,
				},
			})

			if chosen == 0 && recvOk {
				job = value.Interface().(Job)
				ok = true
				found = true
				break
			}
		}

		// 整体阻塞
		if !found {
			chosen, value, recvOk := reflect.Select(allCases)

			// 处理context取消
			if chosen == len(allCases)-1 {
				utils.LogWarn("TaskManager Dispatcher context canceled.")
				// 清理任务
				m.cleanupDebounceMap()
				return
			}

			// 处理channel关闭
			if !recvOk {
				utils.LogInfo("Priority channel %d closed", m.sortedPriorities[chosen])
				return
			}

			// 正常收到任务
			job = value.Interface().(Job)
			ok = true
		}

		// 拿到的 job放入 执行队列
		if ok {
			select {
			case m.execChan <- job:
			case <-m.ctx.Done():
				return
			}
		}
	}
}

// 清理防抖map中的定时器
func (m *DelayedTaskManager) cleanupDebounceMap() {
	m.debounceMap.Range(func(key any, value any) bool {
		if entry, ok := value.(*debounceEntry); ok {
			entry.timer.Stop()
		}
		m.debounceMap.Delete(key)
		return true
	})
}

// 安全提交任务
func (m *DelayedTaskManager) Submit(job Job) {
	priority := job.Priority()
	ch, exists := m.priorityChans[priority]
	if !exists {
		utils.LogError("Unknown priority %d for job: %s", priority, job.Key())
		return
	}

	select {
	case ch <- job:
	default:
		utils.LogWarn("TaskManager channel is full. drop job: %s", job.Key())
	}
}

// 优雅关闭
func (m *DelayedTaskManager) shutdown() {
	utils.LogInfo("shutting down TaskManager...")
	m.cancel() // 发送停止信号
	// 关闭所有优先级channel
	for _, ch := range m.priorityChans {
		close(ch)
	}
	// 关闭执行channel
	close(m.execChan)
}
