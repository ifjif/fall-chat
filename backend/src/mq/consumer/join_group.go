package consumer

import (
	"context"
	"encoding/json"
	"sync"
	"zzc/fall-chat/src/mess"
	"zzc/fall-chat/src/mq"
	"zzc/fall-chat/src/utils"

	"github.com/segmentio/kafka-go"
)

var (
	joinGroupOnce sync.Once
	joinGroup     *JoinGroup
)

type JoinGroup struct{}

func GetJoinGroup() *JoinGroup {
	if joinGroup != nil {
		return joinGroup
	}

	joinGroupOnce.Do(func() {
		joinGroup = &JoinGroup{}
		joinGroup.init()
	})

	return joinGroup
}

func (h *JoinGroup) init() {
	joinGroup := mq.GetDispatcher().GetOrCreateGroup("join-group", "join-group-group")
	joinGroup.RegisterHandler(h.Handle)
	joinGroup.StartInstance()
}

// 从消息队列中监听到 消息
func (h *JoinGroup) Handle(ctx context.Context, message *kafka.Message, commitFunc func()) error {
	var msg mess.JoinGroup
	err := json.Unmarshal(message.Value, &msg)
	if err != nil {
		// 记录错误
		return err
	}
	utils.LogInfo("收到加入群消息：%v\n", msg)
	commitFunc()
	return nil
}
