package data

import (
	"sync"
	"zzc/fall-chat/src/model"

	"gorm.io/gorm"
)

type UpdateConversationJobDao struct {
	BaseDao
}

var (
	updateConversationJobDao     *UpdateConversationJobDao
	updateConversationJobDaoOnce sync.Once
)

func GetUpdateConversationJobDao(db *gorm.DB) *UpdateConversationJobDao {
	updateConversationJobDaoOnce.Do(func() {
		updateConversationJobDao = &UpdateConversationJobDao{}
		updateConversationJobDao.db = db
	})

	return updateConversationJobDao
}

func (dao *UpdateConversationJobDao) UpdateConvLastMessWithTx(tx *gorm.DB, m *model.IMConversation) error {
	return tx.Where("last_msg_seq < ?", m.LastMsgSeq).
		Select("last_msg_id", "last_msg_preview", "last_msg_sender_id", "last_msg_time", "last_msg_seq").
		Updates(m).Error
}

func (dao *UpdateConversationJobDao) UpdateConvMemberMessSeqWithTx(tx *gorm.DB, m *model.IMConversationMember) error {
	return tx.Where("conv_id = ? and user_id = ? and last_read_seq < ? ", m.ConvID, m.UserID, m.LastReadSeq).
		Select("last_read_seq", "last_read_msg_id").
		Updates(m).Error
}
