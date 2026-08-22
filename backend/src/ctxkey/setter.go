package ctxkey

import "github.com/gofiber/fiber/v3"

// 为Locals 设置 UserID值 通常为 int64/string
func SetUserID(c fiber.Ctx, v any) {
	c.Locals(UserID, v)
}

func SetUserPubID(c fiber.Ctx, v any) {
	c.Locals(UserPubID, v)
}

// 为Locals 设置 Username值 通常为 string
func SetUsername(c fiber.Ctx, v any) {
	c.Locals(Username, v)
}
