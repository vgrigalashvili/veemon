package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/vgrigalashvili/veemon/internal/api/rest"
	"github.com/vgrigalashvili/veemon/internal/api/rest/middleware"
	"github.com/vgrigalashvili/veemon/internal/dto"
	"github.com/vgrigalashvili/veemon/internal/helper"
	"github.com/vgrigalashvili/veemon/internal/repository"
	"github.com/vgrigalashvili/veemon/internal/service"
	"github.com/vgrigalashvili/veemon/internal/token"
)

var (
	ErrUniqueMobileComplaint = errors.New("user with given mobile already exists")
	ErrUniqueEmailComplaint  = errors.New("user with given email already exists")
)

type UserHandler struct {
	userService service.UserService
}

func InitializeUserHandler(rh *rest.RestHandler) {

	api := rh.API

	pasetoMaker, err := token.NewPasetoMaker(rh.SEC)
	if err != nil {
		log.Fatalf("[FATAL] error while creating Paseto maker: %v", err)
	}

	userService := service.UserService{
		Token:    pasetoMaker,
		UserRepo: repository.NewUserRepository(rh.DB),
	}

	userHandler := UserHandler{
		userService: userService,
	}

	authMiddleware := middleware.AuthMiddleware(pasetoMaker)

	// public
	api.Post("/user/sign-up", userHandler.signUp)
	api.Post("/user/sign-in", userHandler.signIn)

	// protected
	api.Get("/user/get-by-mobile", authMiddleware, userHandler.getUserByMobile)
	api.Get("/user/get-by-id", authMiddleware, userHandler.getUserByID)
	api.Get("/user/get-by-email", authMiddleware, userHandler.getUserByEmail)
	api.Patch("/user/update-user", authMiddleware, userHandler.updateUser)
	api.Get("/user/get-all-users", authMiddleware, userHandler.getAllUsers)
	api.Get("/user/count", authMiddleware, userHandler.countUsers)
}

func (uh *UserHandler) signUp(ctx *fiber.Ctx) error {
	var credentials dto.UserSignUp
	if err := ctx.BodyParser(&credentials); err != nil {
		log.Printf("[ERROR] invalid request body: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"data":    "invalid request format.",
		})
	}

	validate := validator.New()
	if err := validate.Struct(credentials); err != nil {
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
			"data":    "validation failed.",
		})
	}

	existingUser, err := uh.userService.FindUserByMobile(credentials.Mobile)
	if err == nil && existingUser.Mobile != "" {
		return ctx.Status(http.StatusConflict).JSON(&fiber.Map{
			"success": false,
			"data":    "mobile is already in use.",
		})
	}

	token, err := uh.userService.AddUser(credentials)
	if err != nil {
		if isUniqueConstraintViolation(err) {
			return ctx.Status(http.StatusConflict).JSON(&fiber.Map{
				"success": false,
				"data":    "a user with this mobile already exists.",
			})
		}
		log.Printf("[ERROR] during sign-up: %v", err)
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"data":    "error during the sign-up process.",
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"success": true,
		"data":    token,
	})
}

func (uh *UserHandler) signIn(ctx *fiber.Ctx) error {
	var credentials dto.UserSignIn
	if err := ctx.BodyParser(&credentials); err != nil {
		log.Printf("[ERROR] invalid request body: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"data":    "invalid request format.",
		})
	}

	user, err := uh.userService.FindUserByMobile(credentials.Mobile)
	if err != nil {
		log.Printf("[ERROR] sign-in failed for mobile %s: user not found or database error: %v", credentials.Mobile, err)
		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"success": false,
			"data":    "invalid credentials.",
		})
	}

	if err := helper.CheckPassword(user.Password, credentials.Password); err != nil {
		log.Printf("[ERROR] invalid password attempt for user ID %s", user.ID)
		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"success": false,
			"data":    "invalid credentials.",
		})
	}

	duration := 24 * time.Hour

	token, tPayload, err := uh.userService.Token.CreateToken(user.Mobile, user.Role, duration)
	if err != nil {
		log.Printf("[ERROR] failed to create token for user ID %s: %v", user.ID, err)
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"data":    "error generating token.",
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"success": true,
		"token":   token,
		"payload": tPayload, // Optional: remove in production if unnecessary
	})
}

func (uh *UserHandler) getUserByMobile(ctx *fiber.Ctx) error {
	mobileParam := ctx.Query("mobile")
	log.Printf("[DEBUG] received user mobile parameter: %s", mobileParam)

	user, err := uh.userService.FindUserByMobile(mobileParam)
	if err != nil {
		log.Printf("[ERROR] user not found for mobile: %s: %v", mobileParam, err)
		return ctx.Status(http.StatusNotFound).JSON(&fiber.Map{
			"success": false,
			"data":    "user not found.",
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"success": true,
		"data":    user,
	})
}
func (uh *UserHandler) getUserByID(ctx *fiber.Ctx) error {
	idParam := ctx.Query("id")
	log.Printf("[DEBUG] received user ID parameter: %s", idParam)
	userID, err := uuid.Parse(idParam)
	if err != nil {
		log.Printf("[ERROR] invalid user ID format: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"data":    "invalid user ID format.",
		})
	}

	user, err := uh.userService.FindUserByID(userID)
	if err != nil {
		log.Printf("[ERROR] user not found for ID %s: %v", userID, err)
		return ctx.Status(http.StatusNotFound).JSON(&fiber.Map{
			"success": false,
			"data":    "user not found.",
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"success": true,
		"data":    user,
	})
}
func (uh *UserHandler) getUserByEmail(ctx *fiber.Ctx) error {
	email := ctx.Query("email")
	if email == "" {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"data":    "email query parameter is required",
		})
	}

	user, err := uh.userService.FindUserByEmail(email)
	if err != nil {
		log.Printf("[ERROR] user not found for email %s: %v", email, err)
		return ctx.Status(http.StatusNotFound).JSON(&fiber.Map{
			"success": false,
			"data":    "user not found",
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"success": true,
		"data":    user,
	})
}

func (uh *UserHandler) updateUser(ctx *fiber.Ctx) error {
	// Parse the user ID from query parameters (or include it in the body, if preferred).
	idParam := ctx.Query("id")
	if idParam == "" {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"data":    "user ID query parameter is required.",
		})
	}

	userID, err := uuid.Parse(idParam)
	if err != nil {
		log.Printf("[ERROR] invalid user ID format: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"data":    "invalid user ID format.",
		})
	}

	// Parse the request body into the DTO.
	var updateData dto.UpdateUser
	if err := ctx.BodyParser(&updateData); err != nil {
		log.Printf("[ERROR] invalid request body: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"data":    "invalid request format.",
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
			"data":    "validation failed.",
		})
	}

	// Call the service to update the user.
	updatedUser, err := uh.userService.UpdateUser(userID, updateData)
	if err != nil {
		log.Printf("[ERROR] failed to update user with ID %s: %v", userID, err)
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"data":    "error updating user.",
		})
	}

	// Return the updated user.
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"success": true,
		"data":    updatedUser,
	})
}

func (uh *UserHandler) getAllUsers(ctx *fiber.Ctx) error {
	users, err := uh.userService.GetAllUsers()
	if err != nil {
		log.Printf("[ERROR] error retrieving users: %v", err)
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"data":    "Error retrieving users.",
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"success": true,
		"data":    users,
	})
}

func (uh *UserHandler) countUsers(ctx *fiber.Ctx) error {
	count, err := uh.userService.CountUsers()
	if err != nil {
		log.Printf("[ERROR] error while counting users: %v", err)
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"data":    "error while counting users.",
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"success": true,
		"data":    count,
	})
}

func isUniqueConstraintViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505" // SQLSTATE 23505 indicates unique violation
	}
	return false
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
