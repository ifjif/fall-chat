package service

import (
	"fmt"
	"zzc/fall-chat/src/database"
	"zzc/fall-chat/src/utils"

	"gorm.io/gorm"
)

type BaseService struct {
	name string
	sf   *utils.SnowFlake
	db   *gorm.DB
}

func (bs *BaseService) Service() {}

func (bs *BaseService) Name() string {
	return bs.name
}

func (bs *BaseService) Inject() {
	bs.sf = utils.SF
	bs.db = database.DB
}

func (bs *BaseService) LogWarn(op string, err error, fields ...any) error {
	return fmt.Errorf("%s.%s failed: %w | Context: %v", bs.name, op, err, fields)
	// utils.LogWarn("[%s] %s failed: %v | Context: %v", bs.name, op, err, fields)
}

func (bs *BaseService) LogError(op string, err error, fields ...any) error {
	return fmt.Errorf("%s.%s failed: %w | Context: %v", bs.name, op, err, fields)
	// utils.LogError("[%s] %s failed: %v | Context: %v", bs.name, op, err, fields)
}
