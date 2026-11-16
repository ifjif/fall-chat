package controller

import (
	"zzc/fall-chat/src/errs"

	"github.com/gofiber/fiber/v3"
)

type Base struct {
	name string
}

func (bs *Base) Name() string {
	return bs.name
}

func (bs *Base) SetRoute(f fiber.Router) {}

func (bs *Base) Inject() {}

func (bc *Base) Wrap(op string, err error, fields ...any) error {
	return errs.Wrap(op, err, fields...)
}

func (bc *Base) BindBoyErr(err error) error {
	return bc.Wrap("body bind", err)
}
