package command

type AddGroupRequest struct {
	UserID    int64
	Username  string
	GroupID   int64
	Greeting  string
	InvitedID *int64
	ApplyType int
}

type HandleGroupRequest struct {
	HandleByUserID int64
	Username       string
	HandleID       int64
	Action         int
}
