package job

import "time"

type UpdateConversationJob struct {
	key      string
	duration time.Duration
}

func (uc *UpdateConversationJob) Key() string {
	return uc.key
}

func (uc *UpdateConversationJob) Delay() time.Duration {
	return uc.duration
}

func (uc *UpdateConversationJob) Execute() {
}

func (uc *UpdateConversationJob) OnComplete() {
}
