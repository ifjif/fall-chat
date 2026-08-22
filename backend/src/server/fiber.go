package server

import (
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
			// 业务错误
			if err, ok := err.(*errs.BizError); ok {
				return c.Status(err.Status).JSON(fiber.Map{
					"code":    err.Code,
					"status":  err.Status,
					"message": err.Message,
				})
			}

			// fiber 错误
			if err, ok := err.(*fiber.Error); ok {
				return c.Status(err.Code).JSON(fiber.Map{
					"code":    err.Code,
					"status":  err.Code,
					"message": err.Message,
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
