package model

type UserOrGroupBaseInfoRow struct {
	ID           int64
	PeerNickname string
	Avatar       string
	TypeInfo     string // 'user' / 'group'
}

func (ug *UserOrGroupBaseInfoRow) IsUser() bool {
	return ug.TypeInfo == "user"
}

func (ug *UserOrGroupBaseInfoRow) IsGroup() bool {
	return ug.TypeInfo == "group"
}
