package middleware

import (
	"errors"
	"time"
	"zzc/fall-chat/src/errs"
	"zzc/fall-chat/src/keys/ctxkey"

	"github.com/gofiber/fiber/v3"
)

func LoggerMiddleware(c fiber.Ctx) error {
	start := time.Now()
	err := c.Next()
	latency := time.Since(start)

	status := c.Response().StatusCode()
	logMsg := "HTTP Request"

	reqLogger := ctxkey.GetLogger(c)

	commonFields := []any{
		"method", c.Method(),
		"path", c.Path(),
		"status", status,
		"latency", latency.String(),
	}

	if err != nil {
		commonFields = append(commonFields, "error", err.Error())
		// 携带上下文
		var appErr *errs.AppError
		if errors.As(err, &appErr) {
			commonFields = append(commonFields, appErr.ExtractContext()...)
		}
		reqLogger.Error(logMsg, commonFields...)
	} else {
		reqLogger.Info(logMsg, commonFields...)
	}

	return err
}
