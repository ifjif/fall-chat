package dao

import (
	"context"
	"zzc/fall-chat/src/cursor"
	"zzc/fall-chat/src/data/executor"
	"zzc/fall-chat/src/filter"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/queryparam"
)

type EmojiOwnDao interface {
	CreateBatch(ctx context.Context, m []model.EmojiOwn, batchSize int) error
	List(ctx context.Context, userID int64, qp queryparam.EmojiOwnList, limit int) ([]*model.EmojiPack, error)
}

type emojiOwnDao struct {
	db executor.Executor
}

func NewEmojiOwnDao(db executor.Executor) EmojiOwnDao {
	return &emojiOwnDao{
		db: db,
	}
}

func (dao *emojiOwnDao) CreateBatch(ctx context.Context, m []model.EmojiOwn, batchSize int) error {
	return dao.db.CreateInBatches(m, batchSize).Error
}

func (dao *emojiOwnDao) List(
	ctx context.Context,
	userID int64,
	qp queryparam.EmojiOwnList,
	limit int,
) ([]*model.EmojiPack, error) {
	var result []*model.EmojiPack
	ownPacks := dao.db.Model(&model.EmojiOwn{}).
		Select("distinct pack_id, deleted_at").
		Where("user_id = ?", userID).
		Scopes(filter.BuildFilterScope(qp.Owner.FilterDef, qp.Owner.Filter, nil))

	draftPack := dao.db.Raw("select 1,1")
	// 不是已购买
	if qp.SourceType != 3 {
		// 草稿
		draftPack = dao.db.Model(&model.EmojiPack{}).
			Select("id, deleted_at").
			Where("owner_id = ?", userID).
			Where("status = 1").
			Where("source = 2")
	}

	subquery := dao.db.Raw("? UNION ALL ?", ownPacks, draftPack)

	err := dao.db.Table("(?) AS o", subquery).
		Joins("JOIN emoji_packs p on o.pack_id = p.id").
		Select("p.*").
		Scopes(
			filter.BuildFilterScope(qp.Final.FilterDef, qp.Final.Filter, nil),
			cursor.BuildSortScope(qp.Sort.SortDef, qp.Sort.Sort, nil),
		).
		Limit(limit).
		Find(&result).Error

	return result, err
}
