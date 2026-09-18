package request

import "encoding/json"

type WSMessage struct {
	Type      string          `json:"type"` // 如chat chat_ack join leave system heartbeat
	Timestamp int64           `json:"timestamp"`
	Data      json.RawMessage `json:"data"` // 业务数据，延迟序列化
}

type ChatData struct {
	From        string `json:"from"` // 发送者id
	To          string `json:"to"`   // 接收者id
	Conv        string `json:"conv"` // 会话id
	Content     string `json:"content"`
	ConvType    int    `json:"conv_type"`     // 会话类型 0:单 1:群
	MsgType     int    `json:"msg_type"`      // 消息类型 1:文本 2:图片 3:语言 4:视频
	ClientMsgID string `json:"client_msg_id"` // 本地的临时id
}
