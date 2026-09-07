package model

type UserOrGroupRelationRow struct {
	RelationType  string `gorm:"column:relation_type"` // "user" or "group"
	ID            string `gorm:"column:id"`            // public_id
	Name          string `gorm:"column:name"`          // 用户名或群名
	Avatar        string `gorm:"column:avatar"`
	Nickname      string `gorm:"column:nickname"`       // 仅用户
	Description   string `gorm:"column:description"`    // 仅群
	JoinMode      *int   `gorm:"column:join_mode"`      // 仅群
	FriendStatus  *int   `gorm:"column:friend_status"`  // 仅好友
	RequestStatus *int   `gorm:"column:request_status"` // 仅好友
	GroupUserID   *int64 `gorm:"column:group_user_id"`  // 仅群，在群中则存在
}

/*
'user'

when fs.status = 1 then 'already_friend'
when fsq.status = 0 then 'request_pending'
else 'can_add'

'group'

when gm.user_id is not null then 'already_joined'
when g.join_mode = 2 then 'cannot_join'
when g.join_mode = 1 then 'need_approval'
else 'can_join'
*/
func (u *UserOrGroupRelationRow) GetRelationStatusText() string {
	if u.RelationType == "user" {
		if u.FriendStatus != nil && *u.FriendStatus == 1 {
			return "already_friend"
		}

		if u.RequestStatus != nil && *u.RequestStatus == 0 {
			return "request_pending"
		}

		return "can_add"
	}

	if u.RelationType == "group" {
		if u.GroupUserID != nil {
			return "already_joined"
		}
		if u.RequestStatus != nil && *u.RequestStatus == 0 {
			return "request_pending"
		}

		if *u.JoinMode == 0 {
			return "can_join"
		}

		if *u.JoinMode == 1 {
			return "need_approval"
		}

		if *u.JoinMode == 2 {
			return "cannot_join"
		}

	}

	return "unknown"
}
