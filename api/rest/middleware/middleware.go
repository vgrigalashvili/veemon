package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ResponseDurationLogger(c *fiber.Ctx) error {
	start := time.Now()
	err := c.Next()
	if err != nil {
		return err
	}

	duration := time.Since(start).Milliseconds()

	c.Set("X-Custom-Duration", fmt.Sprintf("%dms", duration))

	return nil
}
