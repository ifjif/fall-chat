package data

import (
	"zzc/fall-chat/src/model"

	"gorm.io/gorm"
)

type ConversationDao struct {
	db *gorm.DB
}

func NewConversationDao(db *gorm.DB) *ConversationDao {
	return &ConversationDao{
		db: db,
	}
}

func (dao *ConversationDao) ListByUserID(userID int64) ([]model.IMConversationRow, error) {
	var rows []model.IMConversationRow
	err := dao.db.Table("im_conversation_members cm").
		Joins("inner join im_conversations cs on cm.conv_id = cs.id").
		Where("cm.user_id = ?", userID).
		Select("cs.*, cm.peer_id,cm.peer_nickname,cm.peer_avatar").
		Find(&rows).Error

	return rows, err
}
