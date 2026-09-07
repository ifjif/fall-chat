package taskmanager

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type TestJob struct {
	id               int
	count            int
	priority         int
	duration         time.Duration
	CompleteCallback func()
}

func NewtTestJob(id, count, priority int, duration time.Duration, CompleteCallback func()) *TestJob {
	return &TestJob{
		id:               id,
		count:            count,
		priority:         priority,
		duration:         duration,
		CompleteCallback: CompleteCallback,
	}
}

func (tj *TestJob) Key() string {
	return fmt.Sprintf("test-%d", tj.id)
}

func (tj *TestJob) Delay() time.Duration {
	return tj.duration
}

func (tj *TestJob) Priority() int {
	return tj.priority
}

func (tj *TestJob) Execute() error {
	fmt.Printf("执行任务: id[%d] priority[%d] count[%d]\n", tj.id, tj.priority, tj.count)
	return nil
}

func (tj *TestJob) RetryCount() int {
	return 0
}

func (tj *TestJob) MaxRetries() int {
	return 0
}

func (tj *TestJob) OnComplete() {
	if tj.CompleteCallback != nil {
		tj.CompleteCallback()
	}
}

func Test1(t *testing.T) {
	var wg sync.WaitGroup
	callback := func() {
		wg.Done()
		fmt.Println("任务完成")
	}

	priority := 0
	for i := range 100 {
		if i > 50 {
			priority = 1
		}
		job := NewtTestJob(i, i, priority, 500*time.Millisecond, callback)
		wg.Add(1)
		TaskManager.Submit(job)
	}

	wg.Wait()
}

type RetryJob struct {
	retryCount       int
	CompleteCallback func()
}

func NewRetryJob() *RetryJob {
	return &RetryJob{}
}

func (rj *RetryJob) Key() string {
	return "key"
}

func (rj *RetryJob) Delay() time.Duration {
	return 0
}

func (rj *RetryJob) RetryCount() int {
	return rj.retryCount
}

func (rj *RetryJob) MaxRetries() int {
	return 3
}

func (rj *RetryJob) Execute() error {
	fmt.Println("执行任务，模拟报错")

	if rj.RetryCount() <= rj.MaxRetries() {
		rj.retryCount++
		return fmt.Errorf("这是模拟报错：准备第 %d 重试", rj.RetryCount())
	}
	return nil
}

func (rj *RetryJob) Priority() int {
	return 1
}

func (rj *RetryJob) OnComplete() {
	if rj.CompleteCallback != nil {
		rj.CompleteCallback()
	}
}

func TestRetry(t *testing.T) {
	time.Sleep(1 * time.Second)
	var wg sync.WaitGroup
	wg.Add(1)

	job := NewRetryJob()
	job.CompleteCallback = func() {
		wg.Done()
	}
	TaskManager.Submit(job)

	wg.Wait()

	/*
		INFO: 正常
		执行任务，模拟报错
		[2026/09/07 19:20:09] [WARN]  Job key failed (attempt 1/3). Retrying in 2s. Error: 这是模拟报错：当前重试:1
		执行任务，模拟报错
		[2026/09/07 19:20:11] [WARN]  Job key failed (attempt 2/3). Retrying in 4s. Error: 这是模拟报错：当前重试:2
		执行任务，模拟报错
		[2026/09/07 19:20:15] [WARN]  Job key failed (attempt 3/3). Retrying in 8s. Error: 这是模拟报错：当前重试:3
		执行任务，模拟报错



		WARN: 重试失败
		执行任务，模拟报错
		[2026/09/07 19:22:54] [WARN]  Job key failed (attempt 1/3). Retrying in 2s. Error: 这是模拟报错：准备第 1 重试
		执行任务，模拟报错
		[2026/09/07 19:22:56] [WARN]  Job key failed (attempt 2/3). Retrying in 4s. Error: 这是模拟报错：准备第 2 重试
		执行任务，模拟报错
		[2026/09/07 19:23:00] [WARN]  Job key failed (attempt 3/3). Retrying in 8s. Error: 这是模拟报错：准备第 3 重试
		执行任务，模拟报错
		[2026/09/07 19:23:08] [WARN]  Job key failed after 3 retries. Last error: 这是模拟报错：准备第 4 重试. Dropping job
	*/
}

var _ = DelayedTaskManager{}
