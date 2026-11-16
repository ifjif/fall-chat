package dao

import (
	"context"
	"errors"
	"zzc/fall-chat/src/data/executor"
	"zzc/fall-chat/src/model"

	"gorm.io/gorm"
)

type ConversationDao interface {
	// 创建
	CreateOne(ctx context.Context, m *model.IMConversation) error

	// 创建一个，返回是否duplicate (true=是)
	CreateOneDup(ctx context.Context, m *model.IMConversation) (bool, error)

	// Select("id, public_id,conv_type").
	FindBriefBySingleKey(ctx context.Context, sk string) (*model.IMConversation, error)

	// 找用户的所有会话
	ListByUserID(ctx context.Context, userID int64) ([]*model.IMConversationRow, error)

	// 找单个会话信息
	SingleByConvID(ctx context.Context, convID, userID int64) (*model.IMConversationRow, error)

	// 通过ids 获取 用户或群的基本信息
	GetUserOrGroupBaseInfoBatch(ctx context.Context, ids []int64) ([]model.UserOrGroupBaseInfoRow, error)

	// 更新会话最后一条消息
	UpdateLastMessage(ctx context.Context, m *model.IMConversation) error
}

type conversationDao struct {
	db executor.Executor
}

func NewConversationDao(db executor.Executor) ConversationDao {
	return &conversationDao{
		db: db,
	}
}

func (dao *conversationDao) FindBriefBySingleKey(ctx context.Context, sk string) (*model.IMConversation, error) {
	var conv model.IMConversation
	err := dao.db.
		Where("single_key = ?", sk).
		Select("id, public_id,conv_type").
		Take(&conv).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return nil, err
	}

	return &conv, err
}

func (dao *conversationDao) CreateOneDup(ctx context.Context, m *model.IMConversation) (bool, error) {
	err := dao.db.Omit("last_msg_time").Create(m).Error

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true, nil
	}

	return false, err
}

func (dao *conversationDao) CreateOne(ctx context.Context, m *model.IMConversation) error {
	return dao.db.Omit("last_msg_time").Create(m).Error
}

func (dao *conversationDao) ListByUserID(ctx context.Context, userID int64) ([]*model.IMConversationRow, error) {
	var rows []*model.IMConversationRow
	err := dao.db.Table("im_conversation_members cm").
		Joins("inner join im_conversations cs on cm.conv_id = cs.id").
		Where("cm.user_id = ?", userID).
		Select("cs.*, cm.peer_id,cm.last_read_seq").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, err
}

func (dao *conversationDao) SingleByConvID(ctx context.Context, convID, userID int64) (*model.IMConversationRow, error) {
	var row model.IMConversationRow
	err := dao.db.Table("im_conversations cs").
		Joins("inner join im_conversation_members cm on cm.conv_id = cs.id and cm.user_id = ?", userID).
		Select("cs.*, cm.peer_id,cm.last_read_seq").
		Where("cs.id = ?", convID).
		Take(&row).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &row, err
}

func (dao *conversationDao) GetUserOrGroupBaseInfoBatch(ctx context.Context, ids []int64) ([]model.UserOrGroupBaseInfoRow, error) {
	var infoResult []model.UserOrGroupBaseInfoRow
	err := dao.db.Raw("? union all ?",
		dao.db.Select("id,nickname as peer_nickname, avatar, 'user' as type_info").
			Model(&model.User{}).
			Where("id in ?", ids),
		dao.db.Select("id, group_name as peer_nickname,avatar, 'group' as type_info").
			Model(&model.GroupInfo{}).
			Where("id in ?", ids),
	).Scan(&infoResult).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return infoResult, nil
}

func (dao *conversationDao) UpdateLastMessage(ctx context.Context, m *model.IMConversation) error {
	return dao.db.Where("last_msg_seq < ?", m.LastMsgSeq).
		Select("last_msg_id", "last_msg_preview", "last_msg_sender_id", "last_msg_time", "last_msg_seq").
		Updates(m).Error
}
