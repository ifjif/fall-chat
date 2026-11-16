package converter

import (
	"zzc/fall-chat/src/command"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/model"
)

func ToRTCCreateRoomCommand(req *request.RTCCreateRoom, userID int64) *command.RTCCreateRoom {
	return &command.RTCCreateRoom{
		RoomType: req.Type,
		UserID:   userID,
	}
}

func ToRTCRoomInviteCommand(req *request.RTCRoomInvite, userID int64) (*command.RTCRoomInvite, error) {
	inviteePubID := req.InviteeID
	inviteeID, err := model.DecodeID(req.InviteeID)
	if err != nil {
		return nil, err
	}

	cmd := &command.RTCRoomInvite{
		UserID:      userID,
		InviteeID:   inviteeID,
		InvitePubID: inviteePubID,
		RoomID:      req.RoomID,
	}
	return cmd, nil
}

func ToRTCRoomInviteAcceptCommand(req *request.RTCRoomInviteAccept, userID int64) (*command.RTCRoomInviteAccept, error) {
	inviterPubID := req.InviterID
	inviterID, err := model.DecodeID(inviterPubID)
	if err != nil {
		return nil, err
	}

	return &command.RTCRoomInviteAccept{
		RoomID:       req.RoomID,
		InviterID:    inviterID,
		InviterPubID: inviterPubID,
		UserID:       userID,
	}, nil
}
