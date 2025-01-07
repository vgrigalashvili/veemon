// Package service provides business logic for user-related operations.
// It interacts with the repository and helper packages for data management and utilities.
package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/vgrigalashvili/veemon/internal/domain"
	"github.com/vgrigalashvili/veemon/internal/dto"
	"github.com/vgrigalashvili/veemon/internal/helper"
	"github.com/vgrigalashvili/veemon/internal/repository"
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
func (us *UserService) AddUser(args domain.User) (string, error) {

	password, err := helper.GeneratePassword()
	log.Printf("[INFO]: Generating password %v", password)
	if err != nil {
		log.Printf("[ERROR] Failed to generate password: %v", err)
		return "", fmt.Errorf("failed to generate password: %w", err)
	}

	expiry := time.Now().AddDate(0, 1, 0)
	hashedPassword, err := helper.HashPassword(password)
	if err != nil {
		log.Printf("[ERROR] Failed to hash password: %v", err)
		return "", fmt.Errorf("failed to hash the password: %w", err)
	}

	user := args

	user.ID = uuid.New()
	user.Password = hashedPassword
	user.Pin = helper.RandomPin()
	user.Role = "backend-developer"
	user.ExpiresAt = &expiry

	log.Printf("[INFO] %v", user)
	log.Printf("[DEBUG] User entity: %+v", user)
	createdUser, err := us.UserRepo.Add(context.Background(), user)
	if err != nil {
		log.Printf("[ERROR - UserService] Failed to add user to the database: %v", err)
		return "", err
	}

	return createdUser.ID.String(), nil
}

// FindUserByID retrieves a user by their unique ID.
// Returns the user or an error if the user could not be found or if an error occurs.
func (us *UserService) GetUserByID(userID uuid.UUID) (domain.User, error) {
	if us.UserRepo == nil {
		log.Printf("[ERROR] UserRepo is not initialized")
		return domain.User{}, fmt.Errorf("UserRepo is not initialized")
	}
	user, err := us.UserRepo.GetByID(context.Background(), userID)
	if err != nil {
		log.Printf("[INFO] not found user by ID %s: %v", userID, err)
		return domain.User{}, fmt.Errorf("could not find user by ID: %w", err)
	}
	return user, nil
}

// FindUserByMobile retrieves a user by their mobile number.
// Return false if the user does not exist, otherwise true.
func (us *UserService) CheckUserByMobile(mobile string) bool {
	log.Printf("[INFO] userRepo : %v", us.UserRepo)
	if us.UserRepo == nil {
		log.Printf("[ERROR] UserRepo is not initialized")
		return false
	}

	if ok := us.UserRepo.CheckByMobile(context.Background(), mobile); !ok {
		log.Printf("[ERROR] User with mobile %s already exists", mobile)
		return false
	}
	return true
}

// FindUserByEmail retrieves a user by their email address.
// Returns the user or an error if the user could not be found or if an error occurs.
func (us *UserService) CheckUserByEmail(email string) (bool, error) {
	if ok := us.UserRepo.CheckByEmail(context.Background(), email); !ok {
		return false, ErrUserNotFound
	}
	return true, nil
}

// UpdateUser updates the details of an existing user based on the provided arguments.
// Returns the updated user or an error if the operation fails.
func (us *UserService) UpdateUser(userID uuid.UUID, arguments dto.UpdateUser) (domain.User, error) {
	// Check if all fields in the arguments are nil
	if arguments.FirstName == nil &&
		arguments.LastName == nil &&
		arguments.Type == nil &&
		arguments.Role == nil &&
		arguments.Mobile == nil &&
		arguments.Email == nil &&
		arguments.Password == nil {
		log.Printf("[INFO] No arguments provided for user update")
		return domain.User{}, fmt.Errorf("nothing to update")
	}
	// Fetch the existing user from the database.
	user, err := us.UserRepo.GetByID(context.Background(), userID)
	if err != nil {
		log.Printf("[ERROR] Error finding user by ID %s: %v", userID, err)
	}

	// Update fields only if they are provided in the arguments.
	if arguments.FirstName != nil {
		user.FirstName = *arguments.FirstName
	}
	if arguments.LastName != nil {
		user.LastName = *arguments.LastName
	}
	if arguments.Type != nil {
		user.Type = *arguments.Type
	}
	if arguments.Role != nil {
		user.Role = *arguments.Role
	}
	if arguments.Mobile != nil {
		user.Mobile = *arguments.Mobile
	}
	if arguments.Email != nil {
		user.Email = *arguments.Email
	}
	if arguments.Password != nil {
		// Hash the new password before updating.
		hashedPassword, err := helper.HashPassword(*arguments.Password)
		if err != nil {
			log.Printf("[ERROR] Failed to hash password: %v", err)
			return domain.User{}, fmt.Errorf("failed to hash the password: %w", err)
		}
		user.Password = hashedPassword
	}
	// if arguments.ExpiresAt != nil {
	// 	user.ExpiresAt = arguments.ExpiresAt
	// }

	// save the updated user to the database.
	updatedUser, err := us.UserRepo.Update(context.Background(), user)
	if err != nil {
		log.Printf("[ERROR] Failed to update user: %v", err)
		return domain.User{}, fmt.Errorf("failed to update user: %w", err)
	}

	return updatedUser, nil
}

// GetAllUsers retrieves all user records from the database.
// Returns a slice of users or an error if the operation fails.
func (us *UserService) GetAllUsers(limit, offset int) ([]domain.User, error) {
	users, err := us.UserRepo.GetAll(context.Background(), limit, offset)
	if err != nil {
		log.Printf("[ERROR] Error getting all users: %v", err)
		return nil, fmt.Errorf("could not get all users: %w", err)
	}
	return users, nil
}

// GetUserByMobile retrieves a user by their mobile number.
// Returns the user or an error if the user could not be found or if an error occurs.
func (us *UserService) GetUserByMobile(mobile string) (domain.User, error) {
	user, err := us.UserRepo.GetByMobile(context.Background(), mobile)
	if err != nil {
		log.Printf("[ERROR] Error finding user by mobile %s: %v", mobile, err)
		return domain.User{}, fmt.Errorf("could not find user by mobile: %w", err)
	}
	return user, nil
}

// GetUserById retrieves a user by their unique ID.
// Returns the user or an error if the user could not be found or if an error occurs.
func (us *UserService) GetUserById(userID uuid.UUID) (domain.User, error) {
	user, err := us.UserRepo.GetByID(context.Background(), userID)
	if err != nil {
		log.Printf("[ERROR] Error finding user by ID %s: %v", userID, err)
		return domain.User{}, fmt.Errorf("could not find user by ID: %w", err)
	}
	return user, nil
}
