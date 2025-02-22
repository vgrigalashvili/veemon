package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vgrigalashvili/veemon/api/rest"
	"github.com/vgrigalashvili/veemon/internal/dto"
	"github.com/vgrigalashvili/veemon/internal/repository"
	"github.com/vgrigalashvili/veemon/internal/service"
	"github.com/vgrigalashvili/veemon/pkg/helper"
	"github.com/vgrigalashvili/veemon/pkg/validator"
)

type AuthHandler struct {
	validator   *validator.CustomValidator
	authService *service.AuthService
}

func InitializeAuthHandler(rh *rest.RestHandler) {

	authRequests := rh.API.Group("api/auth")

	validator := validator.NewValidator()
	userRepository := repository.NewUserRepository(rh.Querier)
	userService := service.NewUserService(userRepository)
	authService := service.NewAuthService(rh.Token, userService)
	authHandler := &AuthHandler{
		authService: authService,
		validator:   validator,
	}

	// public
	authRequests.Post("/sign-up", authHandler.signUp)
	// api.Post("/sign-in", authHandler.signIn)

}

func (ah *AuthHandler) signUp(ctx *fiber.Ctx) error {
	var request dto.AuthSignUp
	clientIP := ctx.IP()
	log.Printf("client ip: %v", clientIP)
	userAgent := ctx.Get("User-Agent")
	log.Printf("user agent: %v", userAgent)
	deviceType := "Desktop"
	if strings.Contains(strings.ToLower(userAgent), "mobile") {
		deviceType = "Mobile"
	}
	log.Printf("device type: %v", deviceType)
	// Parse Request Body
	if err := ctx.BodyParser(&request); err != nil {
		log.Printf("[ERROR] invalid request body: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	if request.Email != "" {
		request.Email = helper.NormalizeEmail(request.Email)
	}
	// Validate Request Struct
	if err := ah.validator.ValidateStruct(&request); err != nil {
		log.Printf("[ERROR] validation failed: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Validation error: " + err.Error(),
		})
	}
	result, err := ah.authService.HandleSignUpProcesses(ctx, request)
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

// func (uh *AuthHandler) signIn(ctx *fiber.Ctx) error {
// 	var credentials dto.AuthSignIn
// 	if err := ctx.BodyParser(&credentials); err != nil {
// 		log.Printf("[ERROR] invalid request body: %v", err)
// 		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
// 			"success": false,
// 			"data":    ErrInvalidRequestJSON.Error(),
// 		})
// 	}
// 	token, err := uh.authService.SignIn(credentials)
// 	if err != nil {

// 		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
// 			"success": false,
// 			"data":    "invalid credentials.",
// 		})
// 	}

// 	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
// 		"success": true,
// 		"data":    token,
// 	})
// }
