package request

type AddGroup struct {
	GroupID   string `json:"target_id" validate:"required"`
	Greeting  string `json:"greeting" validate:"max=255"`
	InvitedID string `json:"invited_id,omitempty"`
	ApplyType int    `json:"apply_type" validate:"oneof=0 1 2"`
}

type ListGroup struct {
	Status  []int  `json:"status"`
	Next    string `json:"next"`
	SortKey string `json:"sort_key"`
}

type HandleGroup struct {
	ID     string `json:"id" validate:"required"`
	Action int    `json:"action" validate:"oneof=1 2"` // 1通过，2:拒绝
}
