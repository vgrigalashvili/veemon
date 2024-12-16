package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/vgrigalashvili/veemon/internal/dto"
	"github.com/vgrigalashvili/veemon/internal/helper"
	"github.com/vgrigalashvili/veemon/internal/token"
)

type AuthService struct {
	Token       token.Maker
	UserService UserService
}

// func (as *AuthService) SignUp(args dto.UserSignUp) (string, error) {
// 	log.Printf("[DEBUG] Starting sign-up process for mobile: %s", args.Mobile)

// 	_, err := as.UserService.FindUserByMobile(args.Mobile)
// 	if err != nil {
// 		log.Printf("[ERROR!] Error checking existing user: %v", err)
// 		return "", err
// 	}

//		log.Printf("[DEBUG] Creating new user: %+v", args.Mobile)
//		token, err := as.UserService.AddUser(args)
//		if err != nil {
//			log.Printf("[ERROR - authService ] Failed to add user: %v", err)
//			return "", err
//		}
//		log.Printf("[DEBUG] User created successfully: %s", token)
//		return token, nil
//	}
func (as *AuthService) SignUp(args dto.UserSignUp) (string, error) {
	log.Printf("[DEBUG] Starting sign-up process for mobile: %s", args.Mobile)

	existingUser, err := as.UserService.FindUserByMobile(args.Mobile)
	if err != nil {
		log.Printf("[ERROR] Error checking existing user: %v", err)
		return "", err
	}

	if existingUser != nil && existingUser.ID != uuid.Nil {
		return "", fmt.Errorf("user with mobile %s already exists", args.Mobile)
	}

	log.Printf("[DEBUG] Creating new user with mobile: %s", args.Mobile)

	token, err := as.UserService.AddUser(args)
	if err != nil {
		log.Printf("[ERROR - AuthService] Failed to add user: %v", err)
		return "", err
	}

	log.Printf("[DEBUG] User created successfully: %s", token)
	return token, nil
}

func (as *AuthService) SignIn(args dto.UserSignIn) (*token.Payload, error) {

	existedUser, err := as.UserService.FindUserByMobile(args.Mobile)
	if err != nil {
		log.Printf("[ERROR] sign-in failed for mobile %s: user not found or database error: %v", args.Mobile, err)
		return &token.Payload{}, err
	}

	if err := helper.CheckPassword(existedUser.Password, args.Password); err != nil {
		log.Printf("[ERROR] invalid password attempt for user ID %s", existedUser.ID)
		return &token.Payload{}, err
	}

	duration := 24 * time.Hour

	tokenPayload, err := as.Token.CreateToken(existedUser.ID, existedUser.Email, existedUser.Role, duration)
	if err != nil {
		log.Printf("[ERROR] failed to create token for user ID %s: %v", existedUser.ID, err)
		return &token.Payload{}, err
	}

	return tokenPayload, nil
}

func (as *AuthService) ForgotPassword() (string, error) {
	return "", nil
}

func (as *AuthService) CreateVerifyEmail(ctx context.Context, email, secretCode string) (string, error) {
	return "", nil
}
