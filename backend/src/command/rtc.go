package command

type RTCCreateRoom struct {
	RoomType string
	UserID   int64
}

type RTCRoomInvite struct {
	UserID      int64
	RoomID      string
	InviteeID   int64
	InvitePubID string
}

type RTCRoomInviteAccept struct {
	UserID       int64
	RoomID       string
	InviterID    int64
	InviterPubID string
}
