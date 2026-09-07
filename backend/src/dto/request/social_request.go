package request

type SocialAddFriendRequest struct {
	TargetID  string `json:"target_id" validate:"required"`
	Greeting  string `json:"greeting" validate:"max=255"`
	ApplyType int    `json:"apply_type" validate:"oneof=0 1 2"`
}
type SocialAddGroupRequestRequest struct {
	GroupID   string `json:"target_id" validate:"required"`
	Greeting  string `json:"greeting" validate:"max=255"`
	InvitedID string `json:"invited_id,omitempty"`
	ApplyType int    `json:"apply_type" validate:"oneof=0 1 2"`
}

type SocialSearchFriendRequestRequest struct {
	FilterStatus []int `json:"filter_status"`
}

type SocialSearchGroupRequestRequest struct {
	FilterStatus []int `json:"filter_status"`
}

type SocialHandleFriendRequestRequest struct {
	ID     string `json:"id" validate:"required"`
	Action int    `json:"action" validate:"oneof=1 2"` // 1:通过，2:拒绝
}

type SocialHandleGroupRequestRequest struct {
	ID     string `json:"id" validate:"required"`
	Action int    `json:"action" validate:"oneof=1 2"` // 1通过，2:拒绝
}
