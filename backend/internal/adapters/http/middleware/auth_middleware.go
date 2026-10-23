package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/services"
	"github.com/morng-dev/erp/pkg/utils"
	"github.com/redis/go-redis/v9"
)

type AuthMiddleware struct {
	jwtSecret          string
	clientRedis        *redis.Client
	permissionsService services.PermissionsService
}

func NewAuthMiddleware(jwtSecret string, clientRedis *redis.Client, permissionsService services.PermissionsService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret:          jwtSecret,
		clientRedis:        clientRedis,
		permissionsService: permissionsService,
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
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(entities.ErrorResponse{
				Message: "invalid user ID in token",
			},
			)
		}
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
		allowed, err := m.permissionsService.HasPermission(c.UserContext(), userID, permission)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(entities.ErrorResponse{
				Message: "error checking permissions",
			})
		}
		if !allowed {
			return c.Status(fiber.StatusForbidden).JSON(entities.ErrorResponse{
				Message: "permission denined",
			})
		}
		return c.Next()
	}
}

func (m *AuthMiddleware) Rolerequire(allowRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(entities.ApiResponse{
				Success: false,
				Message: "ไม่พบข้อมูล role",
			})
		}
		for _, allowRole := range allowRoles {
			if role == allowRole {
				return c.Next()
			}
		}
		return c.Status(fiber.StatusForbidden).JSON(entities.ApiResponse{
			Success: false,
			Message: "ไม่มีสิทธิ์เข้าถึง",
		})
	}

}
