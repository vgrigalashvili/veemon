package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/vgrigalashvili/veemon/api/rest"
	"github.com/vgrigalashvili/veemon/api/rest/middleware"
	"github.com/vgrigalashvili/veemon/internal/domain"
	"github.com/vgrigalashvili/veemon/internal/dto"
	"github.com/vgrigalashvili/veemon/internal/repository"
	"github.com/vgrigalashvili/veemon/internal/service"
)

type UserHandler struct {
	userService *service.UserService
	handleError func(ctx *fiber.Ctx, err error) error // error handler function for handling API errors.
}

func InitializeUserHandler(rh *rest.RestHandler) {

	api := rh.API
	errorHandler := &rest.DefaultAPIErrorHandler{}
	userRepository := repository.NewUserRepository(rh.Querier)
	userService := service.NewUserService(userRepository)

	userHandler := &UserHandler{
		userService: userService,
		handleError: errorHandler.HandleError,
	}

	authMiddleware := middleware.AuthMiddleware(rh.Token)

	// protected
	api.Post("/user/add", authMiddleware, userHandler.add)
	// api.Get("/user/get", authMiddleware, userHandler.get)
	// api.Patch("/user/update", authMiddleware, userHandler.update)
}

// @Summary Add a User
// @Description Creates a new user in the system.
// @Tags Users
// @Accept json
// @Produce json
// @Param user body dto.CreateUser true "User Data"
// @Success 201 {object} dto.StandardResponse
// @Failure 400 {object} dto.StandardResponse
// @Failure 500 {object} dto.StandardResponse
// @Router /user/add [post]
func (uh *UserHandler) add(ctx *fiber.Ctx) error {

	// Parse the request body into the DTO.
	var userData dto.CreateUser
	if err := ctx.BodyParser(&userData); err != nil {
		log.Printf("[ERROR] invalid request body: %v", err)
		return uh.handleError(ctx, rest.ErrInvalidRequestJSON)
	}

	validate := validator.New()
	if err := validate.Struct(userData); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			validationMessages := buildValidationErrorMessages(validationErrors)
			return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"success": false,
				"data":    validationMessages,
			})
		}
		log.Printf("[ERROR] validation error: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"data":    rest.ErrValidationField,
		})
	}

	user := domain.User{
		FirstName: userData.FirstName,
		LastName:  userData.LastName,

		Email: userData.Email,
		Role:  userData.Role,
	}

	userID, err := uh.userService.Create(user)
	if err != nil {
		log.Printf("[ERROR] failed to create user: %v", err)
		if errors.Is(err, ErrUniqueMobileComplaint) {
			return ctx.Status(http.StatusConflict).JSON(&fiber.Map{
				"success": false,
				"data":    ErrUniqueMobileComplaint,
			})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"data":    "failed to create user.",
		})
	}

	return ctx.Status(http.StatusCreated).JSON(&fiber.Map{
		"success": true,
		"data":    userID,
	})
}

func buildValidationErrorMessages(validationErrors validator.ValidationErrors) []string {
	var validationMessages []string
	fieldNames := map[string]string{
		"Email":    "email",
		"Password": "password",
	}
	for _, validationErr := range validationErrors {
		fieldName := fieldNames[validationErr.Field()]
		if fieldName == "" {
			fieldName = validationErr.Field()
		}
		switch validationErr.Tag() {
		case "required":
			validationMessages = append(validationMessages, fmt.Sprintf("validation failed: `%s` is required", fieldName))
		case "email":
			validationMessages = append(validationMessages, fmt.Sprintf("validation failed: `%s` is required.", fieldName))
		case "min":
			validationMessages = append(validationMessages, fmt.Sprintf("validation failed: `%s` must be at least %s characters long.", fieldName, validationErr.Param()))
		default:
			validationMessages = append(validationMessages, fmt.Sprintf("validation failed: `%s` field is invalid.", fieldName))
		}
	}
	return validationMessages
}
