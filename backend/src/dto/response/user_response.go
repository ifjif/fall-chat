package response

import (
	"time"
	"zzc/fall-chat/src/model"
)

type LoginUserResponse struct {
	AccessToken string       `json:"access_token"`
	ExpiresIn   int64        `json:"expires_in"`
	User        UserResponse `json:"user"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginResult struct {
	User             *model.User
	RefreshToken     string
	RefreshExpiresIn int64
	AccessToken      string
	AccessExpiresIn  int64
}

type SearchDataResult struct {
	Type           string `json:"type"` // "user" or "group"
	ID             string `json:"id"`   // public_id
	Name           string `json:"name"` // 用户名或群名
	Avatar         string `json:"avatar"`
	Nickname       string `json:"nickname,omitempty"`    // 仅用户
	Description    string `json:"description,omitempty"` // 仅群
	JoinMode       *int   `json:"join_mode,omitempty"`   // 仅群
	RelationStatus string `json:"relation_status"`
}

type SearchResult struct {
	Result     []*SearchDataResult `json:"result"`
	NextCursor string              `json:"next_cursor"`
}

type UserProfileResponse struct {
	ID         string    `json:"id"`
	Avatar     string    `json:"avatar"`
	Username   string    `json:"username"`
	Nickname   string    `json:"nickname"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	CreatedAt  time.Time `json:"created_at"`
	Status     int       `json:"status"`
	SearchFlag int       `json:"search_flag"`
}
