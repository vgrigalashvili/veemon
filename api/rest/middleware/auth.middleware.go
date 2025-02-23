package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vgrigalashvili/veemon/pkg/token"
)

var (
	ErrAuthHeaderRequired               = errors.New("authorization header required")
	ErrInvalidAuthorizationHeaderFormat = errors.New("authorization header format is not valid")
	ErrInvalidOrExpiredToken            = errors.New("invalid or expired token")
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func AuthMiddleware(tm token.Maker) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		authHeader := ctx.Get(authorizationHeaderKey)
		if authHeader == "" {
			log.Println("[WARN] Missing authorization header")

			return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
				"success": false,
				"data":    ErrAuthHeaderRequired.Error(),
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != authorizationTypeBearer {
			log.Printf("[WARN] Invalid authorization header format: %s", authHeader)

			return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
				"success": false,
				"data":    ErrInvalidAuthorizationHeaderFormat.Error(),
			})
		}

		tokenString := parts[1]

		payload, err := tm.VerifyToken(tokenString)
		if err != nil {
			log.Printf("[ERROR] Token verification failed: %v", err)

			return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
				"success": false,
				"data":    ErrInvalidOrExpiredToken.Error(),
			})
		}

		ctx.Locals(authorizationPayloadKey, payload)
		ctx.Locals("userID", payload.UserID)
		ctx.Locals("userRole", payload.Role)

		return ctx.Next()
	}
}
