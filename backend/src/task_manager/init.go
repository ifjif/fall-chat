package taskmanager

import (
	"context"
	"sort"
)

var TaskManager *DelayedTaskManager

func init() {
	initDelayTaskMananger()
}

func initDelayTaskMananger() {
	ctx, cancel := context.WithCancel(context.Background())
	TaskManager = &DelayedTaskManager{
		priorityChans: map[int]chan Job{
			0: make(chan Job, 10000),
			1: make(chan Job, 10000),
		},
		execChan: make(chan Job, 10000),
		ctx:      ctx,
		cancel:   cancel,
	}

	// 初始化排序后的优先级列表
	for p := range TaskManager.priorityChans {
		TaskManager.sortedPriorities = append(TaskManager.sortedPriorities, p)
	}
	sort.Ints(TaskManager.sortedPriorities)

	// 启动执行器
	go TaskManager.StartExecuter(ctx)
	// 启动调度器
	go TaskManager.StartDispatcher()
}
