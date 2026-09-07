package response

import "time"

type IMConversationListResponse struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Avatar          string    `json:"avatar"`
	Type            string    `json:"type"`
	PeerID          string    `json:"peer_id"`
	LastMsgID       int64     `json:"last_msg_id"`
	LastMsgPreview  string    `json:"last_msg_preview"`
	LastMsgSenderID int64     `json:"last_msg_sender_id"`
	LastMsgTime     time.Time `json:"last_msg_time"`
}
