package dao

import (
	"context"
	"errors"
	"zzc/fall-chat/src/data/executor"
	"zzc/fall-chat/src/model"

	"gorm.io/gorm"
)

type GroupRequestDao interface {
	PendingAndAgreeCount(ctx context.Context, groupID, userID int64) (int64, error)

	CreateOne(ctx context.Context, m *model.GroupRequest) error

	List(ctx context.Context, id int64, filterStatus []int, limit int) ([]model.GroupRequestRow, error)

	// Select("to_group_id, from_user_id, status, apply_type").
	FindBriefOne(ctx context.Context, id int64) (*model.GroupRequest, error)

	UpdateOne(ctx context.Context, updates *model.GroupRequest) error
}

type groupRequestDao struct {
	db executor.Executor
}

func NewGroupRequestDao(db executor.Executor) GroupRequestDao {
	return &groupRequestDao{
		db: db,
	}
}

func (dao *groupRequestDao) PendingAndAgreeCount(ctx context.Context, groupID, userID int64) (int64, error) {
	var existsCount int64

	err := dao.db.Model(&model.GroupRequest{}).
		Where("to_group_id = ? and from_user_id = ? and status in (0,1)", groupID, userID).
		Count(&existsCount).Error

	return existsCount, err
}

func (dao *groupRequestDao) CreateOne(ctx context.Context, m *model.GroupRequest) error {
	return dao.db.Create(m).Error
}

// 找我管理的群 和 我请求的群:
func (dao *groupRequestDao) List(ctx context.Context, id int64, filterStatus []int, limit int) ([]model.GroupRequestRow, error) {
	// 找我管理的群
	subquery := dao.db.Model(&model.GroupMember{}).
		Where("user_id = ? and role in ?", id, []int{1, 2}).
		Select("group_id")

	query1 := dao.db.Model(&model.GroupRequest{}).
		Where("to_group_id in (?)", subquery).
		Select("id")

	//  找我的请求
	query2 := dao.db.Model(&model.GroupRequest{}).
		Where("from_user_id = ?", id).
		Select("id")

	if len(filterStatus) > 0 {
		query1.Where("status in ?", filterStatus)
		query2.Where("status in ?", filterStatus)
	}

	var ids []int64
	err := dao.db.Raw("select id from (? union all ?) as combined order by id desc limit ?",
		query1, query2, 20).
		Scan(&ids).Error
	if err != nil {
		return nil, err
	}

	var rows []model.GroupRequestRow

	if len(ids) <= 0 {
		return rows, err
	}

	err = dao.db.Model(&model.GroupRequest{}).
		Preload("GroupInfo", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, public_id, group_name, avatar")
		}).
		Preload("ApplyUser", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, public_id, nickname, avatar")
		}).
		Preload("InviterUser", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, public_id, nickname, avatar")
		}).
		Preload("HandleByUser", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, public_id, nickname, avatar")
		}).
		Where("id in ?", ids).
		Order("created_at desc").
		Limit(limit).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (dao *groupRequestDao) FindBriefOne(ctx context.Context, id int64) (*model.GroupRequest, error) {
	m := &model.GroupRequest{}
	m.ID = uint64(id)
	err := dao.db.
		Select("to_group_id, from_user_id, status, apply_type").
		Take(m).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return m, err
}

func (dao *groupRequestDao) UpdateOne(ctx context.Context, updates *model.GroupRequest) error {
	return dao.db.Model(updates).Updates(updates).Error
}
