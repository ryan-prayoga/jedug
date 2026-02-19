package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimiter creates a rate limiter middleware
func RateLimiter(maxRequests int, window time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        maxRequests,
		Expiration: window,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Terlalu banyak request. Coba lagi nanti.",
			})
		},
	})
}

// APIRateLimiter for general API routes (100 req/min)
func APIRateLimiter() fiber.Handler {
	return RateLimiter(100, 1*time.Minute)
}

// UploadRateLimiter for upload routes (10 req/min)
func UploadRateLimiter() fiber.Handler {
	return RateLimiter(10, 1*time.Minute)
}

// ReportRateLimiter for creating reports (5 req/min)
func ReportRateLimiter() fiber.Handler {
	return RateLimiter(5, 1*time.Minute)
}
