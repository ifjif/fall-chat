package model

import "gorm.io/plugin/soft_delete"

/*
* id bigint primary key
* public_id  varchar(32) unique key
* created_at datetime
* updated_at datetime
* deleted_at datetime
*
* user_id bigint
* friend_id bigint
* status tinyint default 1 	// 1:正常 0:拉黑
* remark varchar(50)  // 备注
* source varchar(50)	// 来源(search/group_qrcode/phone)
*
* unique key idx_user_frient(user_id,friend_id,deleted_at)
*
 */

type Friendship struct {
	Model
	UserID   int64                 `gorm:"uniqueIndex:idx_user_friend"`
	FriendID int64                 `gorm:"uniqueIndex:idx_user_friend"`
	Status   int8                  `gorm:"default:2"` // 1:正常 0:拉黑
	Remark   string                `gorm:"size:50"`   // 备注
	Source   string                `gorm:"size:50"`   // 来源(search/group_qrcode/phone)
	DeleteAt soft_delete.DeletedAt `gorm:"uniqueIndex:idx_user_friend"`
}
