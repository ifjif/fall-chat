package dao

import (
	"context"
	"sync"
	"zzc/fall-chat/src/data/executor"
	"zzc/fall-chat/src/model"
)

type ConvMembDao interface {
	Create(ctx context.Context, mms []*model.IMConversationMember) error
	// 创建一个群中的成员会话
	CreateOne(ctx context.Context, m *model.IMConversationMember) error

	// 列出当前会话中 所有成员的简短信息
	ListBriefByConvID(ctx context.Context, convID int64) ([]model.IMConvMembBriefRow, error)

	// 通过会话ID和成员ID列表，找出所有成员简短信息
	ListBriefBatchByConvIDAndMemberIDs(ctx context.Context, convID int64, memberIDs []int64) ([]model.IMConvMembBriefRow, error)

	// 更新最后读的消息seq
	UpdateLastReadSeq(ctx context.Context, m *model.IMConversationMember) error
}

type convMembDao struct {
	db executor.Executor
}

func NewConvMembDao(db executor.Executor) ConvMembDao {
	return &convMembDao{
		db: db,
	}
}

var (
	conversationMemberDaoOnce sync.Once
	conversationMemberDao     ConvMembDao
)

func GetConversationMemberDao() ConvMembDao {
	conversationMemberDaoOnce.Do(func() {
		conversationMemberDao = &convMembDao{}
	})

	return conversationMemberDao
}

func (dao *convMembDao) Create(ctx context.Context, mms []*model.IMConversationMember) error {
	return dao.db.Create(mms).Error
}

func (dao *convMembDao) CreateOne(ctx context.Context, m *model.IMConversationMember) error {
	return dao.db.Create(m).Error
}

func (dao *convMembDao) listBrief(convID int64, memberIDs []int64) ([]model.IMConvMembBriefRow, error) {
	var rows []model.IMConvMembBriefRow
	db := dao.db.Table("im_conversation_members icm").
		Joins("inner join users u on icm.user_id = u.id").
		Joins("left join group_members gm on gm.group_id = icm.peer_id and gm.user_id = icm.user_id").
		Where("icm.conv_id = ?", convID)

	if len(memberIDs) != 0 {
		db.Where("icm.user_id in ?", memberIDs)
	}

	err := db.Select("u.id, u.public_id, u.nickname, gm.nickname as in_group_nickname, u.avatar").
		Find(&rows).Error

	return rows, err
}

func (dao *convMembDao) ListBriefByConvID(ctx context.Context, convID int64) ([]model.IMConvMembBriefRow, error) {
	rows, err := dao.listBrief(convID, nil)
	return rows, err
}

func (dao *convMembDao) ListBriefBatchByConvIDAndMemberIDs(ctx context.Context, convID int64, memberIDs []int64) ([]model.IMConvMembBriefRow, error) {
	if len(memberIDs) == 0 {
		return nil, nil
	}

	rows, err := dao.listBrief(convID, memberIDs)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (dao *convMembDao) UpdateLastReadSeq(ctx context.Context, m *model.IMConversationMember) error {
	return dao.db.Where("conv_id = ? and user_id = ? and last_read_seq < ? ", m.ConvID, m.UserID, m.LastReadSeq).
		Select("last_read_seq", "last_read_msg_id").
		Updates(m).Error
}
