package limiter

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	redis "finance-bot/config"
)

func RateLimiterMiddleware(limit int, window time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()

		key := fmt.Sprintf("rate_limit:%s", ip)

		count, err := redis.Client.Get(redis.Ctx, key).Int()
		if err != nil && err.Error() != "redis: nil" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Redis error",
			})
		}

		if count >= limit {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded",
			})
		}

		pipe := redis.Client.TxPipeline()
		pipe.Incr(redis.Ctx, key)
		if count == 0 {
			pipe.Expire(redis.Ctx, key, window)
		}
		_, err = pipe.Exec(redis.Ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Redis pipeline error",
			})
		}

		return c.Next()
	}
}
