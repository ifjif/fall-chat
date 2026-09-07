package model

import "time"

/*
* id bigint primary key
* public_id varchar(32) unique
* created_at datetime
* updated_at datetime
* deleted_at datetime index
*
* conv_id bigint
* sender_id bigint
* msg_type tinyint // 消息类型 1:文本 2:图片 3:语言 4:视频
* content text
* status tinyint // 消息状态 0:正常 1:已撤回 2:已删除
* send_time datetime
* seq bigint default 0 // 会话内单调递增
*
* key idx_conv_id_seq(conv_id, seq)
*
 */

type IMMessage struct {
	Model
	ConvID   int64     `gorm:"column:conv_id;index:idx_conv_id_seq"`
	SenderID int64     `gorm:"column:sender_id"`
	MsgType  int8      `gorm:"column:msg_type"`
	Content  string    `gorm:"column:content"`
	Status   int8      `gorm:"column:status"`
	SendTime time.Time `gorm:"column:send_time"`
	Seq      int64     `gorm:"column:seq;index:idx_conv_id_seq;default:0"`
}
