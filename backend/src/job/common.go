package job

import (
	taskmanager "zzc/fall-chat/src/task_manager"
	"zzc/fall-chat/src/utils"
)

type Common struct{}

func (c *Common) RetryCount() int {
	return 1
}

func (c *Common) MaxRetries() int {
	return 0
}

func (c *Common) Priority() int {
	return 1
}

func (c *Common) ShouldReplace(old taskmanager.Job) bool {
	return true
}

func (cj *Common) LogWarn(op string, err error, fields ...any) {
	utils.LogWarn("[Job] %s failed: %v | Context: %v", op, err, fields)
}

func (cj *Common) OnComplete() {}
