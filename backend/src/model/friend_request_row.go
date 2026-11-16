package model

type FriendshipRequestRow struct {
	FriendshipRequest
	Nickname string `gorm:"column:nickname"`
	Avatar   string `gorm:"column:avatar"`
}
