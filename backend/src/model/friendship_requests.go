package model

/*
* id bigint primary key
* public_id varchar(32) unique key
* created_at datetime
* updated_at datetime
* deleted_at datetime index
*
* from_user_id bigint
* to_user_id bigint
* message varchar(255)
* status tinyint default 0 //0:待处理 1:已通过 2:已拒绝 3:已过期
*
* index idx_to_user_id_status(to_user_id,status)
*
 */

type FriendshipRequest struct {
	Model
	FromUserID int64
	ToUserID   int64  `gorm:"index:idx_to_user_id_status"`
	Message    string `gorm:"size:255"`                              // 申请消息
	Status     int8   `gorm:"index:idx_to_user_id_status;default:0"` // 0:待处理 1:已同意 2:已拒绝 3:已过期
}
