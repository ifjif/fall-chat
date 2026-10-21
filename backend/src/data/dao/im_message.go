package dao

import (
	"context"
	"zzc/fall-chat/src/data/executor"
	"zzc/fall-chat/src/model"
)

type IMMessageDao interface {
	CreateOne(ctx context.Context, m *model.IMMessage) error
}

type imMessageDao struct {
	db executor.Executor
}

func NewIMMessageDao(db executor.Executor) IMMessageDao {
	return &imMessageDao{
		db: db,
	}
}

func (dao imMessageDao) CreateOne(ctx context.Context, m *model.IMMessage) error {
	return dao.db.Create(m).Error
}
