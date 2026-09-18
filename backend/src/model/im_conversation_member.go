package model

/*
* id bigint primary key
* public_id varchar(32) unique
* created_at datetime
* updated_at datetime
* deleted_at datetime index
*
* conv_id bigint // 会话id
* user_id bigint // 用户id
* peer_id bigint // 对方id
*
* unread_count int // 用户的未读数
* top_flag tinyint // 是否置顶：0-否 1-是
* mute_flag tinyint // 是否免打扰: 0-否 1-是
* last_read_seq bigint // 最后读的消息 序列号
*
* unique key idx_conv_user_id_uk(conv_id, user_id)
* key idx_user_id(user_id,peer_id)
*
*
 */

type IMConversationMember struct {
	Model
	ConvID        int64 `gorm:"column:conv_id;uniqueIndex:idx_conv_user_peer_id_uk"`
	UserID        int64 `gorm:"column:user_id;uniqueIndex:idx_conv_user_peer_id_uk;index:idx_user_peer_id"`
	PeerID        int64 `gorm:"column:peer_id;index:idx_user_peer_id"`
	LastReadSeq   int64 `gorm:"column:last_read_seq"`
	TopFlag       int8  `gorm:"column:top_flag"`
	MuteFlag      int8  `gorm:"column:mute_flag"`
	LastReadMsgID int64 `gorm:"column:last_read_msg_id"`
}

func GenerateConvMember(convID, userID, peerID int64) *IMConversationMember {
	return &IMConversationMember{
		ConvID: convID,
		UserID: userID,
		PeerID: peerID,
	}
}
