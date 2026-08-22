package bind

import (
	"zzc/fall-chat/src/errs"

	"github.com/gofiber/fiber/v3"
)

var ErrBind = errs.NewBizError(4000, fiber.StatusBadRequest, "请求参数错误")

func Body[T any](c fiber.Ctx, log func(err error)) (*T, error) {
	var req T

	if err := c.Bind().Body(&req); err != nil {
		log(err)
		return nil, ErrBind
	}

	return &req, nil
}

func Form[T any](c fiber.Ctx, log func(err error)) (*T, error) {
	var req T

	if err := c.Bind().Form(&req); err != nil {
		log(err)
		return nil, ErrBind
	}

	return &req, nil
}

func Query[T any](c fiber.Ctx, log func(err error)) (*T, error) {
	var req T

	if err := c.Bind().Query(&req); err != nil {
		log(err)
		return nil, ErrBind
	}

	return &req, nil
}
