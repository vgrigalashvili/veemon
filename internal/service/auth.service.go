package service

import (
	"context"
)

type AuthService struct {
}

func (as *AuthService) SignIn() (string, error) {
	return "", nil
}

func (as *AuthService) SignUp() (string, error) {
	return "", nil
}

func (as *AuthService) ForgotPassword() (string, error) {
	return "", nil
}

func (as *AuthService) CreateVerifyEmail(ctx context.Context, email, secretCode string) (string, error) {
	return "", nil
}
