package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

func WebSocketMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !websocket.IsWebSocketUpgrade(c) {
			return fiber.ErrUpgradeRequired
		}

		userID, ok := c.Locals("userID").(uuid.UUID)
		if !ok || userID == uuid.Nil {
			return fiber.ErrUnauthorized
		}

		return c.Next()
	}
}
