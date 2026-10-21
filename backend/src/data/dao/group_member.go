package dao

import (
	"context"
	"errors"
	"zzc/fall-chat/src/data/executor"
	"zzc/fall-chat/src/model"

	"gorm.io/gorm"
)

type GroupMemberDao interface {
	CreateOneGroupMember(ctx context.Context, m *model.GroupMember) error

	Exists(ctx context.Context, groupID, userID int64) (int64, error)

	FindRole(ctc context.Context, groupID int64, userID int64) (*model.GroupMember, error)

	CreateOne(ctx context.Context, member *model.GroupMember) error

	FindUserIDs(ctx context.Context, groupID int64) ([]int64, error)
}

type groupMemberDao struct {
	db executor.Executor
}

func NewGroupMemberDao(db executor.Executor) GroupMemberDao {
	return &groupMemberDao{
		db: db,
	}
}

func (dao *groupMemberDao) CreateOneGroupMember(ctx context.Context, m *model.GroupMember) error {
	return dao.db.Create(m).Error
}

func (dao *groupMemberDao) Exists(ctx context.Context, groupID, userID int64) (int64, error) {
	var existsCount int64
	err := dao.db.Model(&model.GroupMember{}).
		Where("group_id = ? and user_id = ?", groupID, userID).
		Count(&existsCount).Error

	return existsCount, err
}

func (dao *groupMemberDao) FindRole(ctc context.Context, groupID int64, userID int64) (*model.GroupMember, error) {
	m := &model.GroupMember{}
	err := dao.db.Where("group_id = ? and user_id = ?", groupID, userID).
		Select("role").
		Take(m).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return m, err
}

func (dao *groupMemberDao) CreateOne(ctx context.Context, member *model.GroupMember) error {
	return dao.db.Create(member).Error
}

func (dao *groupMemberDao) FindUserIDs(ctx context.Context, groupID int64) ([]int64, error) {
	var result []int64
	err := dao.db.Model(model.GroupMember{}).
		Where("group_id = ?", groupID).
		Pluck("user_id", &result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
