package mess

import "zzc/fall-chat/src/model/cached"

type RoomInvite struct {
	RoomID    string             `json:"room_id"`
	RoomType  cached.RtcRoomType `json:"room_type"`
	InviteeID string             `json:"invitee_id"`
	InviterID string             `json:"inviter_id"`
}
