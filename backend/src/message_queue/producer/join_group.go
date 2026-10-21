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
	joinGroupOnce sync.Once
	joinGroup     *JoinGroup
)

type JoinGroup struct {
	writer *kafka.Writer
}

func newJoinGroup() *JoinGroup {
	return &JoinGroup{
		writer: messagequeue.KWriter(),
	}
}

func GetJoinGroup() *JoinGroup {
	if joinGroup != nil {
		return joinGroup
	}

	joinGroupOnce.Do(func() {
		joinGroup = newJoinGroup()
	})

	return joinGroup
}

func (p *JoinGroup) Publish(ctx context.Context, groupID string, m *mess.JoinGroup) error {
	bv, _ := json.Marshal(m)
	go messagequeue.KWriter().WriteMessages(context.Background(), kafka.Message{
		Topic: "join-group",
		Key:   []byte(groupID),
		Value: bv,
	})

	return nil
}
