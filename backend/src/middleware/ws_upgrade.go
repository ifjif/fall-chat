package middleware

import (
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

func WsUpgrade(c fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}

	return fiber.ErrUpgradeRequired
}
