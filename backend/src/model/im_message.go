package model

import (
	"time"
	"zzc/fall-chat/src/utils"
)

/*
* id bigint primary key
* public_id varchar(32) unique
* created_at datetime
* updated_at datetime
* deleted_at datetime index
*
* conv_id bigint
* sender_id bigint
* msg_type tinyint // 消息类型 1:文本 2:图片 3:视频 4:文件 5:语音 6:语音通话 7:视频通话
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

func IsInvalidIMMsgType(msgType int) bool {
	switch msgType {
	case 1, 2, 3, 4, 5, 6, 7:
		return true
	default:
		return false
	}
}

func (m *IMMessage) EncodeConvID() string {
	return EncodeID(m.ConvID)
}

func (m *IMMessage) EncodeSenderID() string {
	return EncodeID(m.SenderID)
}

func (m *IMMessage) GetStatusText() string {
	switch m.Status {
	case 0:
		return "success"
	case 1:
		return "withdraw"
	case 2:
		return "deleted"
	}
	return "unknown"
}

func (m *IMMessage) GetPreview() string {
	switch m.MsgType {
	case 1:
		return utils.TruncateBytes(m.Content, 256)
	case 2:
		return "[图片]"
	case 3:
		return "[视频]"
	case 4:
		return "[文件]"
	case 5:
		return "[语音]"
	case 6:
		return "[语音通话]"
	case 7:
		return "[视频通话]"
	default:
		return "[未知]"
	}
}
