package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func WebSocketMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !websocket.IsWebSocketUpgrade(c) {
			return fiber.ErrUpgradeRequired
		}
		userID := c.Locals("userID")
		if userID == nil {
			return fiber.ErrUnauthorized
		}

		c.Locals("wsUserID", userID)

		return c.Next()
	}
}
