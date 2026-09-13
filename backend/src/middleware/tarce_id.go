package middleware

import (
	"log/slog"
	"zzc/fall-chat/src/ctxkey"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func TraceIDMiddleware(c fiber.Ctx) error {
	// 先从请求头中获取上游传来的tarce-id
	traceID := c.Get("X-Trace-ID")
	if traceID == "" {
		// 没有，生成
		traceID = uuid.New().String()
	}

	// 创建一个绑定了 trace-id的 logger
	reqLogger := slog.With("trace_id", traceID)
	ctxkey.SetLogger(c, reqLogger)
	// 将 Trace-ID 写入响应头，方便前端或调用方排查问题
	c.Set("X-Trace-ID", traceID)
	return c.Next()
}
