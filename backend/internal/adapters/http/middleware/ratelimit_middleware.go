package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type RatelimitMW struct {
	clientRedis *redis.Client
}

func NewRatelimitMW(clientRedis *redis.Client) *RatelimitMW {
	return &RatelimitMW{
		clientRedis: clientRedis,
	}
}

func (m *RatelimitMW) LimitRequest(maxRequests int64, window time.Duration) fiber.Handler {

	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		key := "ratelimit:login:" + c.IP()

		count, err := m.clientRedis.Incr(ctx, key).Result()
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(
				fiber.Map{
					"success": false,
					"message": "Rate limit service unavailable",
				},
			)
		}

		if count == 1 {
			if err := m.clientRedis.Expire(ctx, key, window).Err(); err != nil {
				_ = m.clientRedis.Del(ctx, key).Err()

				return c.Status(fiber.StatusServiceUnavailable).JSON(
					fiber.Map{
						"success": false,
						"message": "Rate limit service unavailable",
					},
				)
			}
		}

		if count > maxRequests {
			return c.Status(fiber.StatusTooManyRequests).JSON(
				fiber.Map{
					"success": false,
					"message": "เข้าสู่ระบบบ่อยเกินไป กรุณาลองใหม่ภายหลัง",
				},
			)
		}

		return c.Next()
	}
}
