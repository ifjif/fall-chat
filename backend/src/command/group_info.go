package command

type CreateGroup struct {
	CreatorPublicID string
	CreatorID       int64
	Username        string
	Name            string
	Description     string
	Avatar          string
	MaxMembers      int
	SearchFlag      int
	JoinMode        int
}
