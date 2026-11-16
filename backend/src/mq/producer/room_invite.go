package producer

import (
	"context"
	"encoding/json"
	"sync"
	"zzc/fall-chat/src/mess"
	"zzc/fall-chat/src/mq"

	"github.com/segmentio/kafka-go"
)

var (
	roomInvite     *RoomInvite
	roomInviteOnce sync.Once
)

type RoomInvite struct {
	writer *kafka.Writer
}

func newRoomInvite() *RoomInvite {
	return &RoomInvite{
		writer: mq.KWriter(),
	}
}

func GetRoomInvite() *RoomInvite {
	if roomInvite != nil {
		return roomInvite
	}

	roomInviteOnce.Do(func() {
		roomInvite = newRoomInvite()
	})

	return roomInvite
}

func (p *RoomInvite) Publish(ctx context.Context, roomID string, mes *mess.RoomInvite) error {
	v, err := json.Marshal(mes)
	if err != nil {
		return err
	}
	p.writer.WriteMessages(ctx, kafka.Message{
		Topic: "rtc-room-invite",
		Key:   []byte(roomID),
		Value: v,
	})

	return nil
}
