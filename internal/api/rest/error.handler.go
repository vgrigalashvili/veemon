package rest

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type APIErrorHandler interface {
	HandleError(c *fiber.Ctx, err error) error
	LogError(err error)
}

// APIError represents a structured error for API responses.
type APIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

// Error implements the error interface for APIError.
func (e *APIError) Error() string {
	return e.Message
}

// NewAPIError creates a new instance of APIError.
func NewAPIError(statusCode int, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
	}
}

type DefaultAPIErrorHandler struct{}

func (h *DefaultAPIErrorHandler) HandleError(c *fiber.Ctx, err error) error {
	// Log error for debugging
	h.LogError(err)

	// Default response
	statusCode := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if apiErr, ok := err.(*APIError); ok {
		statusCode = apiErr.StatusCode
		message = apiErr.Message
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"error": message,
	})
}

func (h *DefaultAPIErrorHandler) LogError(err error) {
	log.Printf("[ERROR] %v", err)
}
