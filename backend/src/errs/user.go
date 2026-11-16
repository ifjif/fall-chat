package errs

import "github.com/gofiber/fiber/v3"

var (
	UserErrUsernameOrPassword = NewBizError(UserStartCode, fiber.StatusNotFound, "用户名或密码错误")
	UserErrNotExistUser       = NewBizError(UserStartCode+1, fiber.StatusNotFound, "用户不存在")
	UserErrExistUser          = NewBizError(UserStartCode+2, fiber.StatusNotFound, "用户已存在")
	UserErrEmail              = NewBizError(UserStartCode+3, fiber.StatusNotFound, "邮箱错误")
	UserErrAccountDisabled    = NewBizError(UserStartCode+4, fiber.StatusBadRequest, "账号已禁用")
)
