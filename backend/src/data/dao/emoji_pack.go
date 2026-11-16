package dao

import (
	"context"
	"zzc/fall-chat/src/cursor"
	"zzc/fall-chat/src/data/executor"
	"zzc/fall-chat/src/filter"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/queryparam"
)

type EmojiPackDao interface {
	CreateOne(ctx context.Context, m *model.EmojiPack) error
	FindOne(ctx context.Context, id int64) (*model.EmojiPack, error)
	List(ctx context.Context, qb queryparam.EmojiPackList, limit int) ([]*model.EmojiPack, error)
	ExistsStatus(ctx context.Context, id, userID uint64, status int) (bool, error)
	UpdateStatus(ctx context.Context, id uint64, status int) error
}

type emojiPackDao struct {
	db executor.Executor
}

func NewEmojiPackDao(db executor.Executor) EmojiPackDao {
	return &emojiPackDao{
		db: db,
	}
}

func (dao *emojiPackDao) CreateOne(ctx context.Context, m *model.EmojiPack) error {
	return dao.db.Create(m).Error
}

func (dao *emojiPackDao) FindOne(ctx context.Context, id int64) (*model.EmojiPack, error) {
	result := model.EmojiPack{}
	result.ID = uint64(id)
	return &result, dao.db.Take(&result).Error
}

func (dao *emojiPackDao) List(ctx context.Context, qp queryparam.EmojiPackList, limit int) ([]*model.EmojiPack, error) {
	var result []*model.EmojiPack
	qb := dao.db.
		Limit(limit).
		Scopes(
			filter.BuildFilterScope(qp.Filter.FilterDef, qp.Filter.Filter, nil),
			cursor.BuildSortScope(qp.Sort.SortDef, qp.Sort.Sort, nil),
		)

	err := qb.Find(&result).Error

	return result, err
}

func (dao *emojiPackDao) ExistsStatus(ctx context.Context, id, userID uint64, status int) (bool, error) {
	var count int64
	err := dao.db.Model(&model.EmojiPack{}).
		Where("id = ? and owner_id = ? and status = ?", id, userID, status).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

func (dao *emojiPackDao) UpdateStatus(ctx context.Context, id uint64, status int) error {
	u := model.EmojiPack{}
	u.ID = id
	return dao.db.Model(u).Update("status", status).Error
}
