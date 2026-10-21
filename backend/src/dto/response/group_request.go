package response

import "time"

type ListGroupRequestResult struct {
	ID                 string    `json:"id"`
	GroupID            string    `json:"group_id"`
	GroupName          string    `json:"group_name"`
	GroupAvatar        string    `json:"group_avatar"`
	FromUserID         string    `json:"from_user_id"`
	FromUserName       string    `json:"from_user_name"`
	FromUserAvatar     string    `json:"from_user_avatar"`
	InviteUserID       string    `json:"invite_user_id"`
	InviteUserName     string    `json:"invite_user_name"`
	InviteUserAvatar   string    `json:"invite_user_avatar"`
	HandleByUserID     string    `json:"handle_by_user_id"`
	HandleByUserName   string    `json:"handle_by_user_name"`
	HandleByUserAvatar string    `json:"handle_by_user_avatar"`
	Message            string    `json:"message"`
	Status             string    `json:"status"`
	Sender             bool      `json:"sender"`
	ApplyType          string    `json:"apply_type"`
	CreatedAt          time.Time `json:"created_at"`
}
