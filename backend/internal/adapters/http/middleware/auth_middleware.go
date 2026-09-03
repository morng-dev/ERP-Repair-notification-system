package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/pkg/utils"
	"github.com/redis/go-redis/v9"
)

type AuthMiddleware struct {
	jwtSecret   string
	clientRedis *redis.Client
}

func NewAuthMiddleware(jwtSecret string, clientRedis *redis.Client) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret:   jwtSecret,
		clientRedis: clientRedis,
	}
}

func (m *AuthMiddleware) AuthRequire() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Cookies("access_token")
		if tokenString == "" {
			authHeader := c.Get("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(entities.ErrorResponse{
				Message: "missing access token",
			})
		}

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(entities.ErrorResponse{
				Message: "invalid or expired token",
			})
		}
		userID, err := uuid.Parse(claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("userID", userID)
		c.Locals("role", claims.Role)
		return c.Next()
	}
}

func (m *AuthMiddleware) PermissionRequire(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("userID").(uuid.UUID)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(entities.ErrorResponse{
				Message: "invalid user",
			})
		}
		key := "user:permissions:" + userID.String()
		exists, err := m.clientRedis.SIsMember(c.Context(), key, permission).Result()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(entities.ErrorResponse{
				Message: "error checking permissions",
			})
		}
		if !exists {
			return c.Status(fiber.StatusForbidden).JSON(entities.ErrorResponse{
				Message: "ไม่มีสิทธิ์เข้าถึง",
			})
		}
		return c.Next()
	}
}
