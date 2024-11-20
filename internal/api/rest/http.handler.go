// Package rest provides the handlers and middleware for the REST API.
// It includes the definition of the RestHandler struct, which is used to manage dependencies across API routes.
package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vgrigalashvili/veemon/internal/token"
	"gorm.io/gorm"
)

// Common error messages for REST handlers.
var (
	ErrInvalidRequestJSON = "error: invalid JSON body in request" // Error returned when a request body has invalid JSON.
)

// RestHandler is the central structure for handling API routes and their dependencies.
type RestHandler struct {
	API   *fiber.App  // Fiber app instance used for routing.
	DB    *gorm.DB    // Database connection instance.
	Token token.Maker // Token maker instance for authentication and authorization.
	SEC   string      // Symmetric key used for secure operations.
}
