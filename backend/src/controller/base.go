package controller

import (
	"zzc/fall-chat/src/errs"

	"github.com/gofiber/fiber/v3"
)

type BaseController struct {
	name string
}

func (bs *BaseController) Name() string {
	return bs.name
}

func (bs *BaseController) SetRoute(f fiber.Router) {}

func (bs *BaseController) Inject() {}

func (bc *BaseController) Wrap(op string, err error, fields ...any) error {
	return errs.Wrap(op, err, fields...)
}
