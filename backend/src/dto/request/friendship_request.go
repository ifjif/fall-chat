package request

type AddFriend struct {
	TargetID  string `json:"target_id" validate:"required"`
	Greeting  string `json:"greeting" validate:"max=255"`
	ApplyType int    `json:"apply_type" validate:"oneof=0 1 2"`
}

type ListFriend struct {
	Status  []int  `json:"status"`
	Next    string `json:"next"`
	SortKey string `json:"sort_key"`
}

type HandleFriend struct {
	ID     string `json:"id" validate:"required"`
	Action int    `json:"action" validate:"oneof=1 2"` // 1:通过，2:拒绝
}
