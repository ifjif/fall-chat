package command

type CreateGroupCommand struct {
	CreatorID   int64
	Username    string
	Name        string
	Description string
	Avatar      string
	MaxMembers  int
	SearchFlag  int
	JoinMode    int
}
