package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"

// RequestID generates a unique request ID for each incoming request,
// stores it in Fiber locals, and sets it as a response header.
func RequestID() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Use existing request ID from header if present, otherwise generate one.
		requestID := c.Get(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Store in locals for downstream use (e.g., logging).
		c.Locals("requestId", requestID)

		// Set response header.
		c.Set(RequestIDHeader, requestID)

		return c.Next()
	}
}
