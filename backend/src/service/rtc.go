package service

import (
	"context"
	"zzc/fall-chat/src/command"
	"zzc/fall-chat/src/data"
	"zzc/fall-chat/src/data/store"
	"zzc/fall-chat/src/errs"
	"zzc/fall-chat/src/mess"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/model/cached"
	"zzc/fall-chat/src/mq/producer"
)

type RTC struct {
	BaseService
	store              store.RTCStore
	roomInviteProducer *producer.RoomInvite
}

func NewRTC(name string) *RTC {
	bs := BaseService{name: name}
	return &RTC{
		BaseService: bs,
	}
}

func (rs *RTC) Inject() {
	rs.store = data.GetBaseDao()
	rs.roomInviteProducer = producer.GetRoomInvite()
}

func (rs *RTC) CreateRoom(ctx context.Context, cmd *command.RTCCreateRoom) (string, error) {
	userID := cmd.UserID

	roomID, err := rs.store.
		RTCCache().
		GetPeerRoom(ctx, userID)
	// 存在旧房间
	if err == nil && roomID != "" {
		state := rs.store.
			RTCCache().
			GetRoomState(ctx, roomID)

		if state == string(cached.RtcRoomStateActive) || state == string(cached.RtcRoomStateWaiting) {
			return "", errs.RTCErrPeerRoomExists
		}

		// 清除旧的
		rs.store.RTCCache().DelPeerRoom(ctx, userID)
	}

	roomID, err = rs.store.
		RTCCache().
		CreateRoom(ctx, userID, cached.RtcRoomTypeOne2One, 2)
	if err != nil {
		return "", rs.Wrap("CreateRoom", err, "user_id", userID)
	}

	return roomID, nil
}

func (rs *RTC) RoomInvite(ctx context.Context, cmd *command.RTCRoomInvite) error {
	roomID := cmd.RoomID

	room, err := rs.store.
		RTCCache().
		GetRoom(ctx, roomID)
	if err != nil {
		return rs.Wrap("find room", err, "room_id", roomID)
	}

	// 找被邀请者
	rID, err := rs.store.RTCCache().GetPeerRoom(ctx, cmd.InviteeID)
	if err != nil {
		return rs.Wrap("find invitee", err, "room_id", roomID, "invitee_id", cmd.InviteeID)
	}
	if rID != "" {
		return errs.RTCErrPeerRoomInviteeExists
	}

	if room == nil {
		return errs.RTCErrPeerRoomExists
	}

	if cmd.UserID == cmd.InviteeID {
		return errs.RTCErrPeerRoomInviteSelf
	}

	if room.CreatorID != cmd.UserID {
		return errs.RTCErrPeerRoomCreator
	}

	if !room.IsPending() || room.IsExpire() {
		return errs.RTCErrPeerRoomNotExists
	}

	err = rs.store.RTCCache().UpdateRoomState(ctx, roomID, cached.RtcRoomStateWaiting)
	if err != nil {
		return rs.Wrap("update room state", err, "room_id", roomID)
	}

	err = rs.store.RTCCache().SetPeerRoom(ctx, cmd.InviteeID, roomID)
	if err != nil {
		return rs.Wrap("set peer room", err, "invitee_id", cmd.InviteeID, "room_id", roomID)
	}

	// 发送kafka
	mes := mess.RoomInvite{
		RoomID:    roomID,
		RoomType:  room.Type,
		InviteeID: cmd.InvitePubID,
		InviterID: model.EncodeID(cmd.UserID),
	}

	rs.roomInviteProducer.Publish(ctx, roomID, &mes)
	return nil
}

func (rs *RTC) InviteAccept(ctx context.Context, cmd *command.RTCRoomInviteAccept) error {
	userID := cmd.UserID
	inviterID := cmd.InviterID
	roomID := cmd.RoomID

	room, err := rs.store.
		RTCCache().
		GetRoom(ctx, roomID)
	if err != nil {
		return rs.Wrap("find room", err, "room_id", roomID)
	}

	if room == nil || room.CreatorID != inviterID || !room.IsWaiting() {
		return errs.RTCErrPeerRoomNotExists
	}

	rID, err := rs.store.
		RTCCache().
		GetPeerRoom(ctx, userID)
	if err != nil {
		return rs.Wrap("find invitee peer_room", err, "room_id", roomID, "invitee_id", userID)
	}
	if rID == "" {
		return errs.RTCErrPeerRoomNotExists
	}

	err = rs.store.
		RTCCache().
		UpdateRoomState(ctx, roomID, cached.RtcRoomStateActive)
	if err != nil {
		return rs.Wrap("update peer_room to active", err, "room_id", roomID)
	}

	err = rs.store.RTCCache().SetPeerRoom(ctx, userID, roomID)
	if err != nil {
		return rs.Wrap("set peer_room for invitee", err, "room_id", roomID, "invitee_id", userID)
	}
	err = rs.store.RTCCache().SetPeerRoom(ctx, inviterID, roomID)
	if err != nil {
		return rs.Wrap("set peer_room for inviter", err, "room_id", roomID, "inviter_id", inviterID)
	}

	return nil
}
