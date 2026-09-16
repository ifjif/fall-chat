package service

import (
	"zzc/fall-chat/src/database"
	"zzc/fall-chat/src/errs"
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

func (bs *BaseService) Wrap(op string, err error, fields ...any) error {
	return errs.Wrap(op, err, fields...)
}
