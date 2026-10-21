package request

type CreateGroup struct {
	Name        string `json:"name" validate:"min=1,max=100"`
	Description string `json:"description" validate:"max=255"`
	Avatar      string `json:"avatar" validate:"max=255"`
	MaxMembers  int    `json:"max_members" validate:"min=2"`
	SearchFlag  int    `json:"search_flag" validate:"oneof=0 1"`
	JoinMode    int    `json:"join_mode" validate:"oneof=0 1 2"`
}

type GroupInfoGroups struct {
	Filter []int `json:"filter"`
}
