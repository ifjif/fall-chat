package command

type AddFriend struct {
	UserID    int64
	Username  string
	TargetID  int64
	Greeting  string
	ApplyType int
}

type HandleFriendRequest struct {
	UserID   int64
	Username string
	HandleID int64
	Action   int
}
