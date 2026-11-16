package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"zzc/fall-chat/src/delivery"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/mess"
	"zzc/fall-chat/src/mq"
	"zzc/fall-chat/src/utils"

	"github.com/segmentio/kafka-go"
)

var (
	roomInvite     *RoomInvite
	roomInviteOnce sync.Once
)

type RoomInvite struct {
	deliverer delivery.Deliverer
}

func NewRoomInvite(deliverer delivery.Deliverer) *RoomInvite {
	return &RoomInvite{
		deliverer: deliverer,
	}
}

func GetRoomInvite(deliverer delivery.Deliverer) *RoomInvite {
	if roomInvite != nil {
		return roomInvite
	}

	roomInviteOnce.Do(func() {
		roomInvite = NewRoomInvite(deliverer)
		roomInvite.init()
	})

	return roomInvite
}

func (c *RoomInvite) init() {
	group := mq.GetDispatcher().GetOrCreateGroup("rtc-room-invite", "rtc-room-invite-group")
	group.RegisterHandler(c.Handle)
	group.StartInstance()
}

func (c *RoomInvite) Handle(ctx context.Context, msg *kafka.Message, commitFunc func()) error {
	var mes mess.RoomInvite

	if err := json.Unmarshal(msg.Value, &mes); err != nil {
		// dirty data
		commitFunc()
		utils.LogError("json unmarshal error in roominvite consumer handler: error = %v", err)

		return nil
	}

	// pack  message of chat
	inviteChatData := request.InviteData{
		RoomID:    mes.RoomID,
		PeerID:    mes.InviteeID,
		InviterID: mes.InviterID,
	}

	chatData := response.WSMessage{
		Type:      "invite",
		Timestamp: time.Now().UnixMilli(),
		Data:      inviteChatData,
	}

	d, err := json.Marshal(chatData)
	if err != nil {
		commitFunc()
		utils.LogError("json marshal chatData for invitee: error=%v", err)
		return nil
	}

	if err := c.deliverer.Deliver(ctx, inviteChatData.PeerID, d); err != nil {
		return fmt.Errorf("deliver chat message to=%s: %w", inviteChatData.PeerID, err)
	}

	commitFunc()
	return nil
}
