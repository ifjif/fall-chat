package errs

import "github.com/gofiber/fiber/v3"

var (
	UserErrUsernameOrPassword = NewBizError(4000, fiber.StatusNotFound, "用户名或密码错误")
	UserErrNotExistUser       = NewBizError(4001, fiber.StatusNotFound, "用户不存在")
	UserErrExistUser          = NewBizError(4002, fiber.StatusNotFound, "用户已存在")
	UserErrEmail              = NewBizError(4003, fiber.StatusNotFound, "邮箱错误")
	UserErrAccountDisabled    = NewBizError(4004, fiber.StatusBadRequest, "账号已禁用")
)
