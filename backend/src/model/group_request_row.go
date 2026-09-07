package model

type GroupRequestRow struct {
	GroupRequest
	GroupInfo    *GroupInfo `gorm:"foreignKey:ToGroupID"`
	ApplyUser    *User      `gorm:"foreignKey:FromUserID"`
	InviterUser  *User      `gorm:"foreignKey:InviterUserID"`
	HandleByUser *User      `gorm:"foreignKey:HandleByUserID"`
}
