package producer

import (
	"context"
	"encoding/json"
	"sync"
	"zzc/fall-chat/src/mess"
	messagequeue "zzc/fall-chat/src/message_queue"

	"github.com/segmentio/kafka-go"
)

var (
	chatOnce sync.Once
	chat     *Chat
)

type Chat struct {
	writer *kafka.Writer
}

func GetChat() *Chat {
	if chat != nil {
		return chat
	}
	chatOnce.Do(func() {
		chat = &Chat{
			writer: messagequeue.KWriter(),
		}
	})

	return chat
}

func (p *Chat) Publish(ctx context.Context, convID string, to string, v []byte) error {
	msg := mess.Chat{
		Message: v,
		ConvID:  convID,
		To:      to,
	}
	value, _ := json.Marshal(msg)
	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: "chat",
		Key:   []byte(convID),
		Value: value,
	})
}
