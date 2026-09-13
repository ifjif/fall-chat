package middleware

import (
	"time"
	"zzc/fall-chat/src/ctxkey"

	"github.com/gofiber/fiber/v3"
)

func LoggerMiddleware(c fiber.Ctx) error {
	start := time.Now()
	err := c.Next()
	latency := time.Since(start)

	status := c.Response().StatusCode()
	logMsg := "HTTP Request"

	reqLogger := ctxkey.GetLogger(c)

	if err != nil {
		reqLogger.Error(logMsg, "method", c.Method(),
			"path", c.Path(),
			"status", status,
			"latency", latency.String(),
			"error", err,
		)
	} else {
		reqLogger.Info(logMsg, "method", c.Method(),
			"path", c.Path(),
			"status", status,
			"latency", latency.String(),
		)
	}

	return err
}
