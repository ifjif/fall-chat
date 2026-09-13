package model

import "time"

/*
* id bigint unsigned primary key
* public_id varchar(32) unique key
* created_at datetime
* updated_at datetime
* deleted_at datetime index
*
* con_type tinyint 		// 会话类型：0-单聊 1-群聊
* single_key varchar(64) unique key // 关系唯一key
*
* // 最后一条消息
*	last_msg_id bigint
*	last_msg_preview varchar(255)
*	last_msg_sender_id bigint
*	last_msg_time datetime
* last_msg_seq bigint
*
 */

type IMConversation struct {
	Model
	ConvType        int8      `gorm:"column:conv_type"` // 会话类型 0:单聊 1:群聊
	SingleKey       string    `gorm:"size:64;column:single_key;uniqueIndex"`
	LastMsgID       int64     `gorm:"column:last_msg_id"`
	LastMsgPreview  string    `gorm:"column:last_msg_preview;size:255"`
	LastMsgSenderID int64     `gorm:"column:last_msg_sender_id"`
	LastMsgTime     time.Time `gorm:"column:last_msg_time"`
	LastMsgSeq      int64     `gorm:"column:last_msg_seq"`
}

func (imc *IMConversation) GetConvTypeText() string {
	switch imc.ConvType {
	case 0:
		return "peer"
	case 1:
		return "group"
	default:
		return "unknown"
	}
}

func (imc *IMConversation) IsGroupConType() bool {
	return imc.ConvType == 1
}

func (imc *IMConversation) SetLastMessage(lm *IMMessage) {
	imc.LastMsgID = int64(lm.ID)
	imc.LastMsgPreview = lm.Content // TODO:根据类型进行设置
	imc.LastMsgSenderID = lm.SenderID
	imc.LastMsgTime = lm.SendTime
	imc.LastMsgSeq = lm.Seq
}

func GenerateGroupConv(groupID string) *IMConversation {
	return &IMConversation{
		SingleKey: groupID,
		ConvType:  1,
	}
}
