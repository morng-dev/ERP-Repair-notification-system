package middleware

import (
	"fmt"
	"log"
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

func (m *RatelimitMW) RatelimitRequest(maxRequests int64, key string, window time.Duration) fiber.Handler {

	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		keys := fmt.Sprintf("ratelimit:%s:%s", key, c.IP())

		count, err := m.clientRedis.Incr(ctx, keys).Result()
		log.Printf(
			"[RATELIMIT] path=%s key=%s ip=%s",
			c.Path(),
			keys,
			c.IP(),
		)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(
				fiber.Map{
					"success": false,
					"message": "Rate limit service unavailable",
				},
			)
		}
		if count == 1 {
			if err := m.clientRedis.Expire(ctx, keys, window).Err(); err != nil {
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
					"message": "ทำรายการบ่อยเกินไป กรุณาลองใหม่ภายหลัง",
				},
			)
		}

		return c.Next()
	}
}
