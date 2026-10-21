package dao

import (
	"context"
	"errors"
	"zzc/fall-chat/src/data/executor"
	"zzc/fall-chat/src/model"

	"gorm.io/gorm"
)

type EmailDao interface {
	CreateOne(ctx context.Context, pc model.PasswordRestCode) error
	FindByID(ctx context.Context, id int64) (*model.PasswordRestCode, error)
	MarkUsed(ctx context.Context, pc *model.PasswordRestCode) error
}

type emailDao struct {
	db executor.Executor
}

func NewEmailDao(db executor.Executor) EmailDao {
	return &emailDao{db: db}
}

func (dao *emailDao) CreateOne(ctx context.Context, pc model.PasswordRestCode) error {
	return dao.db.Create(pc).Error
}

func (dao *emailDao) FindByID(ctx context.Context, id int64) (*model.PasswordRestCode, error) {
	var prc model.PasswordRestCode
	err := dao.db.Where("id = ?", id).Take(&prc).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &prc, nil
}

func (dao *emailDao) MarkUsed(ctx context.Context, prc *model.PasswordRestCode) error {
	return dao.db.Model(prc).Update("used", true).Error
}
