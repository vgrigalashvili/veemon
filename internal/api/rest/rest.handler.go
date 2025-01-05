package rest

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	db "github.com/vgrigalashvili/veemon/internal/db/sqlc"
	"github.com/vgrigalashvili/veemon/internal/token"
)

// common error messages for REST handlers.
var (
	ErrEmailQueryParamRequired = errors.New("query parameter required: email")
	ErrUnverified              = errors.New("unverified user")              // Error returned when a requester is unverified
	ErrUnauthorized            = errors.New("unauthorized")                 // Error returned when a request is unauthorized
	ErrNotFound                = errors.New("not found")                    // Error returned when a request is not found
	ErrInvalidOrExpiredToken   = errors.New("invalid or expired token")     // Error returned when a request contains an invalid or expired token.
	ErrInvalidMethod           = errors.New("invalid method")               // Error returned when a request method is not supported.
	ErrInvalidQueryParam       = errors.New("invalid query parameter")      // Error returned when a request contains invalid query parameters.
	ErrInvalidRequestJSON      = errors.New("invalid JSON body in request") // Error returned when a request body has invalid JSON.
	ErrValidationField         = errors.New("validation field")             // Error returned when a request	validation field
)

// central structure for handling API routes and their dependencies.
type RestHandler struct {
	API     *fiber.App  // Fiber app instance used for routing.
	Querier *db.Queries // database connection instance.
	Token   token.Maker // token maker instance for authentication and authorization.
	// ErrorHandler APIErrorHandler      // error handler for API requests.
	// SEC string // symmetric key used for secure operations.
}
