package job

import (
	taskmanager "zzc/fall-chat/src/task_manager"
	"zzc/fall-chat/src/utils"
)

type CommonJob struct{}

func (c *CommonJob) RetryCount() int {
	return 1
}

func (c *CommonJob) MaxRetries() int {
	return 0
}

func (c *CommonJob) Priority() int {
	return 1
}

func (c *CommonJob) ShouldReplace(old taskmanager.Job) bool {
	return true
}

func (cj *CommonJob) LogWarn(op string, err error, fields ...any) {
	utils.LogWarn("[Job] %s failed: %v | Context: %v", op, err, fields)
}
