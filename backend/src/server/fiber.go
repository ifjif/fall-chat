package server

import (
	"errors"
	"zzc/fall-chat/src/errs"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var APP *fiber.App

type structValidator struct {
	validate *validator.Validate
}

func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

func InitFiber() {
	APP = fiber.New(fiber.Config{
		BodyLimit:       200 * 1024 * 1024,
		StructValidator: &structValidator{validate: validator.New()},
		ErrorHandler: func(c fiber.Ctx, err error) error {
			var bizErr *errs.BizError
			// 业务错误
			if errors.As(err, &bizErr) {
				return c.Status(bizErr.Status).JSON(fiber.Map{
					"code":    bizErr.Code,
					"status":  bizErr.Status,
					"message": bizErr.Message,
				})
			}

			// fiber 错误
			var fiberErr *fiber.Error
			if errors.As(err, &fiberErr) {
				return c.Status(fiberErr.Code).JSON(fiber.Map{
					"code":    fiberErr.Code,
					"status":  fiberErr.Code,
					"message": fiberErr.Message,
				})
			}

			// 未知错误 500
			return c.Status(500).JSON(fiber.Map{
				"code":    9999,
				"status":  500,
				"message": "服务器内部错误",
			})
		},
	})
}
