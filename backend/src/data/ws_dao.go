package data

import (
	"errors"
	"zzc/fall-chat/src/model"

	"gorm.io/gorm"
)

type ConvCache map[string]*model.IMConversation

type WSDao struct {
	db        *gorm.DB
	convCache ConvCache
}

func NewWSDao(db *gorm.DB) *WSDao {
	return &WSDao{
		db:        db,
		convCache: make(ConvCache),
	}
}

func (dao *WSDao) findConversationBySingleKeyWithDB(db *gorm.DB, singleKey string) (*model.IMConversation, error) {
	// 从缓存中找，没有从数据库
	if conv, ok := dao.convCache[singleKey]; ok {
		return conv, nil
	}

	var conv model.IMConversation
	err := db.
		Where("single_key = ?", singleKey).
		Select("id, public_id,conv_type").
		Take(&conv).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	// 设置缓存
	dao.convCache[singleKey] = &conv

	return &conv, nil
}

func (dao *WSDao) FindConversationBySingleKeyWithTx(tx *gorm.DB, singleKey string) (*model.IMConversation, error) {
	return dao.findConversationBySingleKeyWithDB(tx, singleKey)
}

func (dao *WSDao) FindConversationBySingleKey(singleKey string) (*model.IMConversation, error) {
	return dao.findConversationBySingleKeyWithDB(dao.db, singleKey)
}

func (dao *WSDao) CreateConversationWithTx(tx *gorm.DB, m *model.IMConversation) (bool, error) {
	err := tx.Omit("last_msg_time").Create(m).Error

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true, nil
	}

	return false, err
}

func (dao *WSDao) CreateConversationMemberBatchWithTx(tx *gorm.DB, mms []*model.IMConversationMember) error {
	return tx.Create(mms).Error
}

func (dao *WSDao) UpdateConversation(m *model.IMConversation) error {
	return dao.db.Updates(m).Error
}

func (dao *WSDao) CreateIMMesage(m *model.IMMessage) error {
	return dao.db.Create(m).Error
}

func (dao *WSDao) FindUserProfile(id int64) (*model.User, error) {
	user := &model.User{}
	user.ID = uint64(id)
	err := dao.db.
		Select("id, nickname, avatar").
		Take(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return user, err
}
