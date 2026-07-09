package middleware

import "github.com/gofiber/fiber/v2"

type AuthMiddleware struct {
	jwtSecret string
}

func NewNewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{jwtSecret: jwtSecret}
}

func (m *AuthMiddleware) AuthRequire() fiber.Handler {
	return func(c *fiber.Ctx) error {
		AuthoHeader := c.Get("Authorization")

		if AuthoHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON()
		}
		return c.Next()
	}
}
