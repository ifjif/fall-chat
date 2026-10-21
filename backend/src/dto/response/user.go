package response

import (
	"time"
	"zzc/fall-chat/src/model"
)

type LoginUser struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	User        User   `json:"user"`
}

type User struct {
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

type UserProfile struct {
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
