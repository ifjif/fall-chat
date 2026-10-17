package ctxkey

import (
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

// 从Locals中安全获取 UserID
// 不存在或类型不匹配，返回零值
func GetUserID(c fiber.Ctx) int64 {
	v, ok := c.Locals(UserID).(int64)
	if !ok {
		return 0
	}

	return v
}

func GetUserPubID(c fiber.Ctx) string {
	v, ok := c.Locals(UserPubID).(string)
	if !ok {
		return ""
	}

	return v
}

// 从Locals中安全获取 Username
// 不存在或类型不匹配，返回零值
func GetUsername(c fiber.Ctx) string {
	v, ok := c.Locals(Username).(string)
	if !ok {
		return ""
	}

	return v
}

func GetLogger(ctx context.Context) *slog.Logger {
	var l *slog.Logger

	if ctx != nil {
		l, _ = ctx.Value(loggerKey{}).(*slog.Logger)
	}

	if l != nil {
		return l
	}

	// 没有，返回默认
	return slog.Default()
}
