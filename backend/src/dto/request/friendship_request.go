package request

type AddFriend struct {
	TargetID  string `json:"target_id" validate:"required"`
	Greeting  string `json:"greeting" validate:"max=255"`
	ApplyType int    `json:"apply_type" validate:"oneof=0 1 2"`
}

type ListFriend struct {
	FilterStatus []int `json:"filter_status"`
}

type HandleFriend struct {
	ID     string `json:"id" validate:"required"`
	Action int    `json:"action" validate:"oneof=1 2"` // 1:通过，2:拒绝
}
