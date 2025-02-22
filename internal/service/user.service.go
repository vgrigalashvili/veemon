// Package service provides business logic for user-related operations.
// It interacts with the repository and helper packages for data management and utilities.
package service

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/vgrigalashvili/veemon/internal/domain"
	"github.com/vgrigalashvili/veemon/internal/repository"
	"github.com/vgrigalashvili/veemon/pkg/helper"
)

// UserService is responsible for user-related operations and business logic.
// It combines repository interactions with additional processes such as password hashing and token generation.
type UserService struct {
	UserRepo repository.UserRepository // UserRepository interface for user data access.
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	if userRepo == nil {
		log.Fatalf("[FATAL] UserRepository cannot be nil")
	}
	return &UserService{UserRepo: userRepo}
}
func (us *UserService) Create(args domain.User) (string, error) {
	if us.UserRepo == nil {
		log.Printf("[ERROR] UserRepo is not initialized")
		return "", fmt.Errorf("UserRepo is not initialized")
	}

	password, err := helper.GeneratePassword()
	log.Printf("[INFO]: Generating password %v", password)
	if err != nil {
		log.Printf("[ERROR] Failed to generate password: %v", err)
		return "", fmt.Errorf("failed to generate password: %w", err)
	}

	// expiry := time.Now().AddDate(0, 1, 0)
	hashedPassword, err := helper.HashPassword(password)
	if err != nil {
		log.Printf("[ERROR] Failed to hash password: %v", err)
		return "", fmt.Errorf("failed to hash the password: %w", err)
	}

	user := args

	user.ID = uuid.New()
	user.Password = hashedPassword
	user.Role = "backend-developer"

	log.Printf("[INFO] %v", user)
	log.Printf("[DEBUG] User entity: %+v", user)
	createdUser, err := us.UserRepo.Create(context.Background(), user)
	if err != nil {
		log.Printf("[ERROR - UserService] Failed to add user to the database: %v", err)
		return "", err
	}

	return createdUser.ID.String(), nil
}

// FindUserByID retrieves a user by their unique ID.
// Returns the user or an error if the user could not be found or if an error occurs.
func (us *UserService) GetBID(userID uuid.UUID) (*domain.User, error) {
	if us.UserRepo == nil {
		log.Printf("[ERROR] UserRepo is not initialized")
		return nil, fmt.Errorf("UserRepo is not initialized")
	}
	user, err := us.UserRepo.Read(context.Background(), userID)
	if err != nil {
		log.Printf("[INFO] not found user by ID %s: %v", userID, err)
		return nil, fmt.Errorf("could not find user by ID: %w", err)
	}
	return user, nil
}
