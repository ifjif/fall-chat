package dao

import (
	"context"
	"zzc/fall-chat/src/cursor"
	"zzc/fall-chat/src/data/executor"
	"zzc/fall-chat/src/filter"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/queryparam"
)

type FriendshipRequestDao interface {
	CreateOne(ctx context.Context, fr *model.FriendshipRequest) error

	PendingAndAgreeCount(ctx context.Context, id, targetID int64) (int64, error)

	List(
		ctx context.Context,
		userID uint64,
		qp queryparam.FriendshipRequestList,
		limit int) ([]*model.FriendshipRequestRow, error)

	UpdateStatus(ctx context.Context, targetID int64, updates map[string]any) error

	FindWithToUerID(ctx context.Context, id int64, toUserID int64) (*model.FriendshipRequest, error)
}

type friendshipRequestDao struct {
	db executor.Executor
}

func NewFriendshipRequestDao(db executor.Executor) FriendshipRequestDao {
	return &friendshipRequestDao{
		db: db,
	}
}

func (dao *friendshipRequestDao) CreateOne(ctx context.Context, fr *model.FriendshipRequest) error {
	return dao.db.Omit("pair_key").Create(fr).Error
}

func (dao *friendshipRequestDao) PendingAndAgreeCount(ctx context.Context, id, targetID int64) (int64, error) {
	// 0:待处理 1:已同意
	var existsCount int64
	err := dao.db.
		Model(&model.FriendshipRequest{}).
		Where(`((from_user_id = ? and to_user_id = ?) or (from_user_id = ? and to_user_id = ?)) 
		and status in (0,1)`,
			id, targetID,
			targetID, id,
		).Count(&existsCount).Error

	return existsCount, err
}

func (dao *friendshipRequestDao) List(
	ctx context.Context,
	userID uint64,
	qp queryparam.FriendshipRequestList,
	limit int,
) ([]*model.FriendshipRequestRow, error) {
	mo := []*model.FriendshipRequestRow{}

	fPrefix := func(f string) string {
		return "f." + f
	}

	qc := dao.db.Model(&model.FriendshipRequest{}).
		Table("friendship_requests as f").
		Joins("inner join users u on (u.id = f.from_user_id or u.id = f.to_user_id) and u.id != ?", userID).
		Where("f.to_user_id = ? or f.from_user_id = ?", userID, userID).
		Limit(limit).
		Scopes(
			filter.BuildFilterScope(qp.Filter.FilterDef, qp.Filter.Filter, fPrefix),
			cursor.BuildSortScope(qp.Sort.SortDef, qp.Sort.Sort, fPrefix),
		)

	err := qc.Select("f.*, u.nickname, u.avatar").Find(&mo).Error

	return mo, err
}

func (dao *friendshipRequestDao) UpdateStatus(ctx context.Context, targetID int64, updates map[string]any) error {
	return dao.db.Model(&model.FriendshipRequest{}).
		Where("id = ?", targetID).
		Updates(updates).Error
}

func (dao *friendshipRequestDao) FindWithToUerID(ctx context.Context, id int64, toUserID int64) (*model.FriendshipRequest, error) {
	var row model.FriendshipRequest
	err := dao.db.Where("id = ? and to_user_id = ?", id, toUserID).Take(&row).Error
	return &row, err
}
