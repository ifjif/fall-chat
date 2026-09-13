package data

import (
	"errors"
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
		Select("cs.*, cm.peer_id,cm.peer_nickname,cm.peer_avatar,cm.last_read_seq").
		Find(&rows).Error

	return rows, err
}

func (dao *ConversationDao) SingleByConvID(convID, userID int64) (*model.IMConversationRow, error) {
	var row model.IMConversationRow
	err := dao.db.Table("im_conversations cs").
		Joins("inner join im_conversation_members cm on cm.conv_id = cs.id and cm.user_id = ?", userID).
		Select("cs.*, cm.peer_id,cm.peer_nickname,cm.peer_avatar,cm.last_read_seq").
		Where("cs.id = ?", convID).
		Take(&row).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &row, err
}
