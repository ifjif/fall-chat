package dao

import (
	"context"
	"zzc/fall-chat/src/cursor"
	"zzc/fall-chat/src/data/executor"
	"zzc/fall-chat/src/model"
)

type EmojiItemDao interface {
	CreateBatch(ctx context.Context, m []model.EmojiItem, batch int) error
	Create(ctx context.Context, m []model.EmojiItem) error
	List(ctx context.Context, packID uint64, sort map[string]any, sortDef cursor.SortDefinition, limit int) ([]*model.EmojiItem, error)
}

type emojiItemDao struct {
	db executor.Executor
}

func NewEmojiItemDao(db executor.Executor) EmojiItemDao {
	return &emojiItemDao{
		db: db,
	}
}

func (dao *emojiItemDao) CreateBatch(ctx context.Context, m []model.EmojiItem, batch int) error {
	return dao.db.CreateInBatches(m, batch).Error
}

func (dao *emojiItemDao) Create(ctx context.Context, m []model.EmojiItem) error {
	return dao.db.Create(&m).Error
}

func (dao *emojiItemDao) List(ctx context.Context, packID uint64, sort map[string]any, sortDef cursor.SortDefinition, limit int) ([]*model.EmojiItem, error) {
	var result []*model.EmojiItem

	orderBy := sortDef.OrderBy(nil)
	q := dao.db.Where("pack_id = ?", packID).
		Order(orderBy).
		Limit(limit)

	whereClause, args, err := sortDef.Where(sort, nil)
	if err != nil {
		return nil, err
	}
	if args != nil {
		q.Where(whereClause, args...)
	}

	err = q.Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
