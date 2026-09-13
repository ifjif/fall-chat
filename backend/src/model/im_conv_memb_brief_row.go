package model

type IMConvMembBriefRow struct {
	ID             int64  `gorm:"column:id" json:"-"`
	PublicID       string `gorm:"column:public_id" json:"id"`
	Nickname       string `gorm:"column:nickname" json:"nickname"`
	InGropNickname string `gorm:"column:in_group_nickname" json:"in_group_nickname"`
	Avatar         string `gorm:"column:avatar" json:"avatar"`
}
