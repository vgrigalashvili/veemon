// Package middleware provides custom middleware functions for the Fiber framework.
// This includes utilities like logging response durations.
package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ResponseDurationLogger is a middleware function that logs the duration of each request.
// It calculates the time taken to process the request and sets it in a custom header (`X-Custom-Duration`).
func ResponseDurationLogger(c *fiber.Ctx) error {
	// Record the start time of the request.
	start := time.Now()

	// Process the next middleware or handler in the chain.
	err := c.Next()
	if err != nil {
		// Return the error if one occurs during the request processing.
		return err
	}

	// Calculate the duration of the request in milliseconds.
	duration := time.Since(start).Milliseconds()

	// Add the duration as a custom header in the response.
	c.Set("X-Custom-Duration", fmt.Sprintf("%dms", duration))

	// Return nil to indicate the middleware execution completed successfully.
	return nil
}
