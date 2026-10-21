package dao

import (
	"context"
	"zzc/fall-chat/src/data/executor"
	"zzc/fall-chat/src/model"
)

type MessageSeqDao interface {
	GenerateNextSeq(ctx context.Context, convID int64) (int64, error)
}

type messageSeqDao struct {
	db executor.Executor
}

func NewMessageSeqDao(db executor.Executor) MessageSeqDao {
	return &messageSeqDao{
		db: db,
	}
}

func (dao *messageSeqDao) GenerateNextSeq(ctx context.Context, convID int64) (int64, error) {
	var curSeq int64
	err := dao.db.Model(&model.IMMessage{}).
		Where("conv_id = ?", convID).
		Pluck("COALESCE(MAX(seq), 0)", &curSeq).Error
	if err != nil {
		return 0, err
	}

	newSeq := curSeq + 1

	return newSeq, nil
}
