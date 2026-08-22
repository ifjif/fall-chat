package model

import "gorm.io/plugin/soft_delete"

/*
* id bigint primary key
* public_id varchar(32) unique key
* created_at datetime
* updated_at datetime
* deleted_at datetime
*
* group_id bigint
* user_id bigint
* nickname varchar(100)
* join_source varchar(50) 		// 加入来源 (search/group_qrcode)
* muted tinyint								// 个人禁言: 0:否 1:是
* role tinyint								// 角色: 0:普通 1:管理员 2:群主
*
* unique key idx_group_user_id(group_id, user_id,deleted_at)
* index idx_user_id(user_id)
*
*
 */

type GroupMember struct {
	Model
	GroupID    int64                 `gorm:"uniqueIndex:idx_group_user_id"`
	UserID     int64                 `gorm:"uniqueIndex:idx_group_user_id;index"`
	DeletedAt  soft_delete.DeletedAt `gorm:"uniqueIndex:idx_group_user_id`
	Nickname   string                `gorm:"size:100"`
	JoinSource string                `gorm:"size:50"`   // 加入来源 (search/group_qrcode)
	Muted      int8                  `gorm:"default:0"` // 个人禁言: 0:否 1:是
	Role       int8                  `gorm:"default:0"` // 角色: 0:普通 1:管理员 2:群主
}
