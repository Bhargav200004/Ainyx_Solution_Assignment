package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

// RequestLogger logs each request's method, path, status code, duration,
// and the associated request ID using Uber Zap.
func RequestLogger(log *zap.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		// Process request.
		err := c.Next()

		duration := time.Since(start)
		requestID, _ := c.Locals("requestId").(string)

		log.Info("request completed",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("duration", duration),
			zap.String("requestId", requestID),
		)

		return err
	}
}
