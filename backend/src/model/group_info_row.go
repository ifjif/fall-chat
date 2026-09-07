package model

type GroupInfoRow struct {
	GroupInfo
	MemberCount int `gorm:"column:member_count"`
}

func (gi *GroupInfoRow) OverflowMember(incr int) bool {
	return (gi.MemberCount + incr) > gi.MaxMembers
}
