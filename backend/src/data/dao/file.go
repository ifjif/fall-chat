package dao

import (
	"context"
	"errors"
	"zzc/fall-chat/src/data/executor"
	"zzc/fall-chat/src/model"

	"gorm.io/gorm"
)

type FileDao interface {
	CreateOne(ctx context.Context, fc model.FileConfig) error
	FindOne(ctx context.Context) (*model.FileConfig, error)
}

type fileDao struct {
	db executor.Executor
}

func NewFileDao(db executor.Executor) FileDao {
	return &fileDao{
		db: db,
	}
}

func (dao *fileDao) CreateOne(ctx context.Context, fc model.FileConfig) error {
	return dao.db.Create(&fc).Error
}

func (dao *fileDao) FindOne(ctx context.Context) (*model.FileConfig, error) {
	var fc model.FileConfig
	err := dao.db.Take(&fc).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &fc, nil
}
