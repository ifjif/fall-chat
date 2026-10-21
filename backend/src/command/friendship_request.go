package command

type AddFriend struct {
	ID        int64
	Username  string
	TargetID  string
	Greeting  string
	ApplyType int
}

type HandleFriendRequest struct {
	ID       int64
	Username string
	HandleID string
	Action   int
}
