// Package middleware provides authentication and authorization middleware for the Fiber framework.
// This includes verifying JWT or PASETO tokens and enforcing access control.
package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vgrigalashvili/veemon/internal/token"
)

var (
	ErrAuthHeaderRequired               = errors.New("authorization header required")
	ErrInvalidAuthorizationHeaderFormat = errors.New("authorization header format is not valid")
	ErrInvalidOrExpiredToken            = errors.New("invalid or expired token")
)

const (
	authorizationHeaderKey  = "authorization"         // Header key for the authorization token.
	authorizationTypeBearer = "bearer"                // Expected token type in the authorization header.
	authorizationPayloadKey = "authorization_payload" // Key used to store token payload in context locals.
)

// AuthMiddleware is a middleware function that validates the authorization token in incoming requests.
// It uses the provided token maker to verify the token and stores the payload in the request context.
func AuthMiddleware(tm token.Maker) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Retrieve the authorization header.
		authHeader := ctx.Get(authorizationHeaderKey)
		if authHeader == "" {
			log.Println("[WARN] Missing authorization header")
			// Return a 401 Unauthorized response if the header is missing.
			return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
				"success": false,
				"data":    ErrAuthHeaderRequired,
			})
		}

		// Split the header into parts and validate the format.
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != authorizationTypeBearer {
			log.Printf("[WARN] Invalid authorization header format: %s", authHeader)
			// Return a 401 Unauthorized response for invalid format.
			return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
				"success": false,
				"data":    ErrInvalidAuthorizationHeaderFormat,
			})
		}

		// Extract the token string from the header.
		tokenString := parts[1]

		// Verify the token using the token maker.
		payload, err := tm.VerifyToken(tokenString)
		if err != nil {
			log.Printf("[ERROR] Token verification failed: %v", err)
			// Return a 401 Unauthorized response if the token is invalid or expired.
			return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
				"success": false,
				"data":    ErrInvalidOrExpiredToken,
			})
		}

		// Store the token payload in the request context for later use.
		ctx.Locals(authorizationPayloadKey, payload)
		ctx.Locals("userID", payload.UserID)
		// Call the next middleware or handler in the chain.
		return ctx.Next()
	}
}
