package response

type WSMessage struct {
	Type      string `json:"type"` // 如chat chat_ack join leave system heartbeat
	Timestamp int64  `json:"timestamp"`
	Data      any    `json:"data"` // 业务数据
}

type WSErr struct {
	Error error `json:"error"`
	Data  any   `json:"data"`
}

type ChatACK struct {
	ConvID      string `json:"conv_id"`
	To          string `json:"to"`
	ServerMsgID string `json:"server_msg_id"`
	ClientMsgID string `json:"client_msg_id"`
	MsgSeq      int64  `json:"msg_seq"`
}

type ChatReply struct {
	ConvID   string `json:"conv_id"`
	To       string `json:"to"`
	MsgID    string `json:"msg_id"`
	SenderID string `json:"sender_id"`
	Content  string `json:"content"`
	MsgSeq   int64  `json:"msg_seq"`
}
type ChatErrData struct {
	ConvID string `json:"conv_id"`
	MsgID  string `json:"msg_id"`
}
