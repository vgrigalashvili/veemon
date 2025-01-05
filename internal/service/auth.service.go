package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/vgrigalashvili/veemon/internal/domain"
	"github.com/vgrigalashvili/veemon/internal/dto"
	"github.com/vgrigalashvili/veemon/internal/helper"
	"github.com/vgrigalashvili/veemon/internal/repository"
	"github.com/vgrigalashvili/veemon/internal/token"
)

var (
	ErrUserAlreadyExists = errors.New("user with this mobile already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrExpiredToken      = errors.New("expired token")
	ErrInvalidToken      = errors.New("invalid token")
)

type AuthService struct {
	Token       token.Maker
	UserService *UserService
}

func NewAuthService(token token.Maker, userService *UserService) *AuthService {
	return &AuthService{
		Token:       token,
		UserService: userService,
	}
}
func (as *AuthService) SignUp(args dto.AuthSignUp) (string, error) {
	if as.UserService == nil {
		return "", errors.New("internal server error: UserService is not initialized")
	}

	userExists, err := as.UserService.CheckUserByMobile(args.Mobile)
	if err != nil && !errors.Is(err, repository.ErrNoRows) {
		log.Printf("[ERROR] not found: %v", err)
		return "", err
	}

	if userExists {
		log.Printf("[ERROR] user with mobile %s already exists", args.Mobile)
		return "", ErrUserAlreadyExists
	}

	newUser := domain.User{
		Mobile: args.Mobile,
	}

	userID, err := as.UserService.AddUser(newUser)
	if err != nil {
		log.Printf("[ERROR] failed to add user: %v", err)
		return "", err
	}

	return userID, nil
}

func (as *AuthService) SignIn(args dto.AuthSignIn) (*token.Payload, error) {
	existedUser, err := as.UserService.GetUserByMobile(args.Mobile)
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
