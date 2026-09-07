package model

type FriendshipRow struct {
	ID       string `json:"id" gorm:"column:public_id"`
	Nickname string `json:"nickname" gorm:"column:nickname"`
	Avatar   string `json:"avatar" gorm:"column:avatar"`
	Remark   string `json:"remark" gorm:"column:remark"`
	Source   string `json:"source" gorm:"column:source"`
}
