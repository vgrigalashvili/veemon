package rest

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type APIErrorHandler interface {
	HandleError(c *fiber.Ctx, err error) error
	LogError(err error)
}

type APIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func (e *APIError) Error() string {
	return e.Message
}

func NewAPIError(statusCode int, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
	}
}

type DefaultAPIErrorHandler struct{}

func (h *DefaultAPIErrorHandler) HandleError(c *fiber.Ctx, err error) error {

	h.LogError(err)

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
