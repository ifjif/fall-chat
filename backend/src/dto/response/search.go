package response

type UserOrGroupDataResult struct {
	Type           string `json:"type"` // "user" or "group"
	ID             string `json:"id"`   // public_id
	Name           string `json:"name"` // 用户名或群名
	Avatar         string `json:"avatar"`
	Nickname       string `json:"nickname,omitempty"`    // 仅用户
	Description    string `json:"description,omitempty"` // 仅群
	JoinMode       *int   `json:"join_mode,omitempty"`   // 仅群
	RelationStatus string `json:"relation_status"`
}

type UserOrGroupResult struct {
	Result     []*UserOrGroupDataResult `json:"result"`
	NextCursor string                   `json:"next_cursor"`
}
