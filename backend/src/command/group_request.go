package command

type AddGroupRequest struct {
	ID        int64
	Username  string
	GroupID   string
	Greeting  string
	InvitedID string
	ApplyType int
}

type HandleGroupRequest struct {
	HandleByUserID int64
	Username       string
	HandleID       string
	Action         int
}
