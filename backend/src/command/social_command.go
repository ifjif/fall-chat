package command

type SocialAddFriendCommand struct {
	ID        int64
	Username  string
	TargetID  string
	Greeting  string
	ApplyType int
}

type SocialAddGroupRequestCommand struct {
	ID        int64
	Username  string
	GroupID   string
	Greeting  string
	InvitedID string
	ApplyType int
}

type SocialHandleFriendRequestCommand struct {
	ID       int64
	Username string
	HandleID string
	Action   int
}

type SocialHandleGroupRequestCommand struct {
	HandleByUserID int64
	Username       string
	HandleID       string
	Action         int
}
