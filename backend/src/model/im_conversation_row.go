package model

type IMConversationRow struct {
	IMConversation
	PeerID       int64  `gorm:"column:peer_id"`
	PeerAvatar   string `gorm:"column:peer_avatar"`
	PeerNickname string `gorm:"column:peer_nickname"`
}
