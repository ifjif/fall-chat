package service

import (
	"zzc/fall-chat/src/errs"
)

type BaseService struct {
	name string
}

func (bs *BaseService) Service() {}

func (bs *BaseService) Name() string {
	return bs.name
}

func (bs *BaseService) Inject() {}

func (bs *BaseService) Wrap(op string, err error, fields ...any) error {
	return errs.Wrap(op, err, fields...)
}
