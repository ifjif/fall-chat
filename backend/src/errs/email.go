package errs

import "github.com/gofiber/fiber/v3"

var EmailErrInvalid = NewBizError(EmailStartCode, fiber.StatusBadRequest, "验证码失效")
