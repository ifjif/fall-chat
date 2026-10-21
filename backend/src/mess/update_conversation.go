package mess

import "time"

type UpdateConv struct {
	ConvID      int64     `json:"conv_id"`   // 会话id
	PeerID      int64     `json:"peer_id"`   // 对方id 用户或群
	Seq         int64     `json:"seq"`       // 消息序列
	ConvType    int       `json:"conv_type"` // 会话类型 0:单聊  1:群
	MsgID       int64     `json""msg_id"`
	MsgPreview  string    `json:"msg_preview"`
	MsgSenderID int64     `json:"msg_sender_id"`
	MsgTime     time.Time `json:"msg_time"`
}
