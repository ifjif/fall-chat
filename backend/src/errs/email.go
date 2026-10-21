package errs

import "github.com/gofiber/fiber/v3"

var EmailErrInvalid = NewBizError(4200, fiber.StatusBadRequest, "验证码失效")
