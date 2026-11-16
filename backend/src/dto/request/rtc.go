package request

type RTCCreateRoom struct {
	Type string `json:"type" validate:"required"`
}

type RTCRoomInvite struct {
	RoomID    string `json:"room_id" validate:"required"`
	InviteeID string `json:"invitee_id" validate:"required"`
}

type RTCRoomInviteAccept struct {
	RoomID    string `json:"room_id" validate:"required"`
	InviterID string `json:"inviter_id" validate:"required"`
}
