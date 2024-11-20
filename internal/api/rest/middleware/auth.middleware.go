// Package middleware provides authentication and authorization middleware for the Fiber framework.
// This includes verifying JWT or PASETO tokens and enforcing access control.
package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vgrigalashvili/veemon/internal/token"
)

const (
	authorizationHeaderKey  = "authorization"         // Header key for the authorization token.
	authorizationTypeBearer = "bearer"                // Expected token type in the authorization header.
	authorizationPayloadKey = "authorization_payload" // Key used to store token payload in context locals.
)

// AuthMiddleware is a middleware function that validates the authorization token in incoming requests.
// It uses the provided token maker to verify the token and stores the payload in the request context.
func AuthMiddleware(tokenMaker token.Maker) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Retrieve the authorization header.
		authHeader := ctx.Get(authorizationHeaderKey)
		if authHeader == "" {
			log.Println("[WARN] Missing authorization header")
			// Return a 401 Unauthorized response if the header is missing.
			return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
				"success": false,
				"data":    "authorization header is required",
			})
		}

		// Split the header into parts and validate the format.
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != authorizationTypeBearer {
			log.Printf("[WARN] Invalid authorization header format: %s", authHeader)
			// Return a 401 Unauthorized response for invalid format.
			return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
				"success": false,
				"data":    "invalid authorization header format",
			})
		}

		// Extract the token string from the header.
		tokenString := parts[1]

		// Verify the token using the token maker.
		payload, err := tokenMaker.VerifyToken(tokenString)
		if err != nil {
			log.Printf("[ERROR] Token verification failed: %v", err)
			// Return a 401 Unauthorized response if the token is invalid or expired.
			return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
				"success": false,
				"data":    "invalid or expired token",
			})
		}

		// Store the token payload in the request context for later use.
		ctx.Locals(authorizationPayloadKey, payload)

		// Call the next middleware or handler in the chain.
		return ctx.Next()
	}
}
