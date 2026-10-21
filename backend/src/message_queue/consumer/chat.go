package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"zzc/fall-chat/src/delivery"
	"zzc/fall-chat/src/keys/ctxkey"
	"zzc/fall-chat/src/mess"
	messagequeue "zzc/fall-chat/src/message_queue"
	"zzc/fall-chat/src/utils"

	"github.com/segmentio/kafka-go"
)

var (
	chatOnce sync.Once
	chat     *Chat
)

type Chat struct {
	deliverer delivery.Deliverer
}

func NewChat(d delivery.Deliverer) *Chat {
	return &Chat{deliverer: d}
}

func (h *Chat) init() {
	chatGroup := messagequeue.GetDispatcher().GetOrCreateGroup("chat", "chat-group")
	chatGroup.RegisterHandler(h.Handle)
	chatGroup.StartInstance()
}

func GetChat(d delivery.Deliverer) *Chat {
	if chat != nil {
		return chat
	}

	chatOnce.Do(func() {
		chat = NewChat(d)
		chat.init()
	})

	return chat
}

// 从消息队列中 监听到消息
// 得到消息
// 通过本地找
// 没有找到，从redis找到rpc地址
// 然后 直接发送
func (h *Chat) Handle(ctx context.Context, message *kafka.Message, commitFunc func()) error {
	logger := ctxkey.GetLogger(ctx)
	var chatMess mess.Chat
	if err := json.Unmarshal(message.Value, &chatMess); err != nil {
		// 脏消息，无法投机，记录警告
		commitFunc()
		logger.Error("[json]", "reason", "chat handle", "error", err)
		return nil
	}
	utils.LogInfo("收到消息：to=%s\n", chatMess.To)

	if err := h.deliverer.Deliver(ctx, chatMess.To, chatMess.Message); err != nil {
		return fmt.Errorf("deliver chat message to=%s: %w", chatMess.To, err)
	}

	commitFunc()
	return nil
}
