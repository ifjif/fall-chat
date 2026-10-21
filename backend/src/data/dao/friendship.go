package dao

import (
	"context"
	"zzc/fall-chat/src/data/executor"
	"zzc/fall-chat/src/model"
)

type FriendshipDao interface {
	// 创建一批
	Create(ctx context.Context, list []model.Friendship) error
	// 找用户的好友
	FriendshipList(ctx context.Context, id int64, limit int) ([]model.FriendshipRow, error)

	// 找两者之间的关系
	FindRelationship(ctx context.Context, id int64, targetID int64) ([]model.Friendship, error)
}

type friendshipDao struct {
	db executor.Executor
}

func NewFriendshipDao(db executor.Executor) FriendshipDao {
	return &friendshipDao{
		db: db,
	}
}

func (dao *friendshipDao) Create(ctx context.Context, list []model.Friendship) error {
	return dao.db.Create(list).Error
}

func (dao *friendshipDao) FindRelationship(ctx context.Context, id int64, targetID int64) ([]model.Friendship, error) {
	friendships := []model.Friendship{}
	err := dao.db.
		Where("(user_id = ? and friend_id = ?) or (user_id = ? and friend_id = ?)",
			targetID, id, // 对方，我
			id, targetID, // 我，对方
		).
		Find(&friendships).Error

	return friendships, err
}

func (dao *friendshipDao) FriendshipList(ctx context.Context, id int64, limit int) ([]model.FriendshipRow, error) {
	var rows []model.FriendshipRow
	err := dao.db.
		Table("friendships as f").
		Joins("inner join users u on u.id = f.friend_id").
		Where("f.user_id = ?", id).
		Select("u.public_id, f.Remark, f.Source, u.nickname, u.avatar").
		Limit(limit).
		Find(&rows).Error

	return rows, err
}
