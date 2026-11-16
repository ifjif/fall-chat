package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"zzc/fall-chat/src/mess"
	"zzc/fall-chat/src/mq"

	"github.com/segmentio/kafka-go"
)

var (
	updateConvMemberOnce sync.Once
	updateConvMember     *UpdateConvMember
)

type UpdateConvMember struct {
	writer *kafka.Writer
}

func newUpdateConvMember() *UpdateConvMember {
	return &UpdateConvMember{
		writer: mq.KWriter(),
	}
}

func GetUpdateConvMember() *UpdateConvMember {
	if updateConvMember != nil {
		return updateConvMember
	}

	updateConvMemberOnce.Do(func() {
		updateConvMember = newUpdateConvMember()
	})
	return updateConvMember
}

func (p *UpdateConvMember) Publish(ctx context.Context, convID string, m *mess.UpdateConv) error {
	messData, err := json.Marshal(m)
	if err != nil {
		// 脏消息
		return fmt.Errorf("json.Marshal 更新会话消息数据: %w", err)
	}
	go mq.KWriter().WriteMessages(context.Background(), kafka.Message{
		Topic: "update-conv-member",
		Key:   []byte(convID),
		Value: messData,
	})

	return nil
}
