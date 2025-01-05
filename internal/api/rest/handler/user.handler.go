package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/vgrigalashvili/veemon/internal/api/rest"
	"github.com/vgrigalashvili/veemon/internal/api/rest/middleware"
	"github.com/vgrigalashvili/veemon/internal/domain"
	"github.com/vgrigalashvili/veemon/internal/dto"
	"github.com/vgrigalashvili/veemon/internal/repository"
	"github.com/vgrigalashvili/veemon/internal/service"
)

// UserHandler is responsible for handling user-related operations.
// It combines service interactions with API request/response handling.
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
	api.Get("/user/get", authMiddleware, userHandler.get)
	api.Patch("/user/update", authMiddleware, userHandler.update)
}

func (uh *UserHandler) add(ctx *fiber.Ctx) error {

	// Parse the request body into the DTO.
	var userData dto.CreateUser
	if err := ctx.BodyParser(&userData); err != nil {
		log.Printf("[ERROR] invalid request body: %v", err)
		return uh.handleError(ctx, rest.ErrInvalidRequestJSON)
	}

	// Validate the input data.
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
		Mobile:    userData.Mobile,
		Email:     userData.Email,
		Role:      userData.Role,
	}
	// Call the service to create the user.
	userID, err := uh.userService.AddUser(user)
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

	// Return the created user.
	return ctx.Status(http.StatusCreated).JSON(&fiber.Map{
		"success": true,
		"data":    userID,
	})
}

func (uh *UserHandler) get(ctx *fiber.Ctx) error {
	if uh.userService == nil {
		log.Printf("[ERROR] user service not initialized")
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"data":    "internal server error: user service not initialized",
		})
	}
	idParam := ctx.Query("id")
	if idParam == "" {
		log.Printf("[INFO] no user ID provided in the request query parameter")
		mobileParam := ctx.Query("mobile")
		if mobileParam == "" {
			log.Printf("[INFO] no user mobile provided in the request query parameter")
			return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"success": false,
				"data":    ErrInvalidQueryParam,
			})
		}
		user, err := uh.userService.GetUserByMobile(mobileParam)
		if err != nil {
			log.Printf("[ERROR] user not found for mobile: %s: %v", mobileParam, err)
			return ctx.Status(http.StatusNotFound).JSON(&fiber.Map{
				"success": false,
				"data":    ErrNotFound,
			})
		}
		return ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"success": true,
			"data":    user,
		})
	}
	// Parse the user ID from query parameters (or include it in the body, if preferred).
	userID, err := uuid.Parse(idParam)
	if err != nil {
		log.Printf("[ERROR] invalid user ID format: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"data":    ErrInvalidUserIDFormat,
		})
	}
	user, err := uh.userService.GetUserByID(userID)
	if err != nil {
		log.Printf("[ERROR] user not found for ID %s: %v", userID, err)
		return ctx.Status(http.StatusNotFound).JSON(&fiber.Map{
			"success": false,
			"data":    ErrNotFound,
		})
	}
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"success": true,
		"data":    user,
	})
}

func (uh *UserHandler) update(ctx *fiber.Ctx) error {
	// Extract userID from the token via middleware
	requesterID := ctx.Locals("userID")
	if requesterID == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"data":    ErrUnauthorized,
		})
	}
	requesterRole := ctx.Locals("userRole")
	if requesterRole == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"data":    ErrUnverified,
		})
	}
	// Parse the user ID from the request URL
	requestedID := ctx.Query("id")
	if requestedID == "" {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"data":    ErrInvalidQueryParam,
		})
	}

	if fmt.Sprint(requestedID) != fmt.Sprint(requesterID) && fmt.Sprint(requesterRole) != "super" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"data":    "you are not allowed to update this user's information",
		})
	}

	// Parse the user ID from query parameters (or include it in the body, if preferred).
	userID, err := uuid.Parse(requestedID)
	if err != nil {
		log.Printf("[ERROR] invalid user ID format: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"data":    ErrInvalidUserIDFormat,
		})
	}

	// Parse the request body into the DTO.
	var updateData dto.UpdateUser
	if err := ctx.BodyParser(&updateData); err != nil {
		log.Printf("[ERROR] invalid JSON body in request: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"data":    ErrInvalidRequestJSON,
		})
	}

	// Validate the update data.
	validate := validator.New()
	if err := validate.Struct(updateData); err != nil {
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
			"data":    ErrValidationField,
		})
	}

	// Call the service to update the user.
	updatedUser, err := uh.userService.UpdateUser(userID, updateData)
	if err != nil {
		log.Printf("[ERROR] failed to update user with ID %s: %v", userID, err)
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			// TODO: update error
			"data": err.Error(),
		})
	}

	// Return the updated user.
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"success": true,
		"data":    updatedUser,
	})
}

func buildValidationErrorMessages(validationErrors validator.ValidationErrors) []string {
	var validationMessages []string
	fieldNames := map[string]string{
		"Mobile":   "mobile",
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
		case "mobile":
			validationMessages = append(validationMessages, fmt.Sprintf("validation failed: `%s` is an invalid mobile.", fieldName))
		case "min":
			validationMessages = append(validationMessages, fmt.Sprintf("validation failed: `%s` must be at least %s characters long.", fieldName, validationErr.Param()))
		default:
			validationMessages = append(validationMessages, fmt.Sprintf("validation failed: `%s` field is invalid.", fieldName))
		}
	}
	return validationMessages
}
