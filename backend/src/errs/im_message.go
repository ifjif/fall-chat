package errs

import "github.com/gofiber/fiber/v3"

var ImMessErrSeqNotMatch = NewBizError(ImMessageStartCode, fiber.StatusBadRequest, "消息seq不匹配")
