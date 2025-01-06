package handler

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/vgrigalashvili/veemon/internal/api/rest"
	"github.com/vgrigalashvili/veemon/internal/dto"
	"github.com/vgrigalashvili/veemon/internal/repository"
	"github.com/vgrigalashvili/veemon/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func InitializeAuthHandler(rh *rest.RestHandler) {

	api := rh.API
	userRepository := repository.NewUserRepository(rh.Querier)
	userService := service.NewUserService(userRepository)
	authService := service.NewAuthService(rh.Token, userService)
	authHandler := &AuthHandler{
		authService: authService,
	}
	// public
	api.Post("/sign-up", authHandler.signUp)
	api.Post("/sign-in", authHandler.signIn)

}

func (ah *AuthHandler) signUp(ctx *fiber.Ctx) error {
	var request dto.AuthSignUp
	if err := ctx.BodyParser(&request); err != nil {
		log.Printf("[ERROR] invalid request body: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"data":    ErrInvalidRequestJSON,
		})
	}

	// Validation logic
	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			validationMessages := buildValidationErrorMessages(validationErrors)
			return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"success": false,
				"data":    validationMessages,
			})
		}
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"data":    ErrValidationField,
		})
	}

	// Check if the user already exists with the given mobile number.
	result, err := ah.authService.HandleSignUpProcesses(request)
	if err != nil {
		log.Printf("[ERROR] failed to sign up: %v", err)
		if err.Error() == "user with this mobile already exists" {
			return ctx.Status(http.StatusConflict).JSON(&fiber.Map{
				"success": false,
				"data":    "you already have registered with this mobile",
			})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"data":    "something went wrong:" + err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (uh *AuthHandler) signIn(ctx *fiber.Ctx) error {
	var credentials dto.AuthSignIn
	if err := ctx.BodyParser(&credentials); err != nil {
		log.Printf("[ERROR] invalid request body: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"data":    ErrInvalidRequestJSON.Error(),
		})
	}
	token, err := uh.authService.SignIn(credentials)
	if err != nil {

		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"success": false,
			"data":    "invalid credentials.",
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"success": true,
		"data":    token,
	})
}
