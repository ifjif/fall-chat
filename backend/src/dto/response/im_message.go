package response

import "time"

type IMMessageResult struct {
	ID       string    `json:"id"`
	ConvID   string    `json:"conv_id"`
	SenderID string    `json:"sender_id"`
	MsgType  int8      `json:"msg_type"`
	Content  string    `json:"content"`
	Status   string    `json:"status"`
	SendTime time.Time `json:"send_time"`
	Seq      int64     `json:"seq"`
}
