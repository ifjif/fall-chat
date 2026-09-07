package model

type SearchFriendRequestRow struct {
	FriendshipRequest
	Nickname string `gorm:"column:nickname"`
	Avatar   string `gorm:"column:avatar"`
}
