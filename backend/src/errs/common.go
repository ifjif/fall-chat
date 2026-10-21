package errs

import "github.com/gofiber/fiber/v3"

var (
	CommonErrJWTTokenGeneration = NewBizError(2000, fiber.StatusInternalServerError, "JWT token生成失败")
	CommonErrJWTTokenExpired    = NewBizError(2001, fiber.StatusUnauthorized, "凭证已过期")
	CommonErrJWTTokenType       = NewBizError(2002, fiber.StatusUnauthorized, "token类型不正确")
)
