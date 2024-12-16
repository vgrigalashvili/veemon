package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/vgrigalashvili/veemon/internal/api/rest"
	"github.com/vgrigalashvili/veemon/internal/dto"
	"github.com/vgrigalashvili/veemon/internal/repository"
	"github.com/vgrigalashvili/veemon/internal/service"
	"github.com/vgrigalashvili/veemon/internal/token"
)

type AuthHandler struct {
	authService service.AuthService
}

func InitializeAuthHandler(rh *rest.RestHandler) {

	api := rh.API

	pasetoMaker, err := token.NewPasetoMaker(rh.SEC)
	if err != nil {
		log.Fatalf("[FATAL] error while creating Paseto maker: %v", err)
	}

	authService := service.AuthService{
		Token: pasetoMaker,
		UserService: service.UserService{
			UserRepo: repository.NewUserRepository(rh.DB),
		},
	}

	authHandler := AuthHandler{
		authService: authService,
	}
	// public
	api.Post("/sign-up", authHandler.signUp)
	api.Post("/sign-in", authHandler.signIn)

}

func (ah *AuthHandler) signUp(ctx *fiber.Ctx) error {
	var credentials dto.UserSignUp
	if err := ctx.BodyParser(&credentials); err != nil {
		log.Printf("[ERROR] invalid request body: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"data":    errInvalidRequestFormat,
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
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"data":    errValidationField,
		})
	}

	result, err := ah.authService.SignUp(credentials)
	if err != nil {
		// Check if the error message contains the duplicate mobile message
		if err.Error() == fmt.Sprintf("user with mobile %s already exists", credentials.Mobile) {
			return ctx.Status(http.StatusConflict).JSON(&fiber.Map{
				"success": false,
				"data":    "a user with this mobile already exists.",
			})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"data":    "something went wrong.",
		})
	}
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (uh *AuthHandler) signIn(ctx *fiber.Ctx) error {
	var credentials dto.UserSignIn
	if err := ctx.BodyParser(&credentials); err != nil {
		log.Printf("[ERROR] invalid request body: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"data":    errInvalidRequestFormat,
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
