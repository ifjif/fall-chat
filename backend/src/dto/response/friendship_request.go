package response

import "time"

type ListFriendRequestResult struct {
	ID         string    `json:"id"`
	FromUserID string    `json:"from_user_id"`
	ToUserID   string    `json:"to_user_id"`
	Sender     bool      `json:"sender"`
	Message    string    `json:"message"`
	Status     string    `json:"status"`
	Avatar     string    `json:"avatar"`
	Nickname   string    `json:"nickname"`
	ApplyType  string    `json:"apply_type"`
	CreatedAt  time.Time `json:"created_at"`
}
