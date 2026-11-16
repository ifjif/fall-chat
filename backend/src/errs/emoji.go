package errs

import "github.com/gofiber/fiber/v3"

var (
	EmojiErrPackNotExists            = NewBizError(EmojiStartCode, fiber.StatusBadRequest, "表情包不存在")
	EmojiErrPublishOutOfBound        = NewBizError(EmojiStartCode+1, fiber.StatusBadRequest, "推送量太多")
	EmojiErrPublishItemCountNotMatch = NewBizError(EmojiStartCode+2, fiber.StatusBadRequest, "推送量不匹配")
)
