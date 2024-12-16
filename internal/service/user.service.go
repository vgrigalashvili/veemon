// Package service provides business logic for user-related operations.
// It interacts with the repository and helper packages for data management and utilities.
package service

import (
	"errors"
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
	// Token    token.Maker               // Token maker for creating and validating tokens.
	UserRepo repository.UserRepository // UserRepository interface for user data access.
}

// AddUser creates a new user based on provided input arguments and saves it to the database.
// Returns the new user's ID as a string or an error if the operation fails.
// func (us *UserService) AddUser(args dto.UserSignUp) (string, error) {
// 	log.Printf("[DEBUG] Inserting user with mobile: %s", args.Mobile)

// 	// Set expiration date for one month from the current time.
// 	expiry := time.Now().AddDate(0, 1, 0)

// 	// Generate a password and handle any errors.
// 	password, err := helper.GeneratePassword()
// 	log.Printf("[INFO]: Generating password %v", password)
// 	if err != nil {
// 		log.Printf("[ERROR] Failed to generate password: %v", err)
// 		return "", fmt.Errorf("failed to generate password: %w", err)
// 	}

// 	// Hash the generated password and handle errors.
// 	hashedPassword, err := helper.HashPassword(password)
// 	if err != nil {
// 		log.Printf("[ERROR] Failed to hash password: %v", err)
// 		return "", fmt.Errorf("failed to hash the password: %w", err)
// 	}

// 	// Construct a new user entity.
// 	user := domain.User{
// 		ID:        uuid.New(), // Generate a new UUID for the user ID.
// 		Mobile:    args.Mobile,
// 		Password:  hashedPassword,
// 		ExpiresAt: &expiry, // Set expiration date.
// 	}

// 	// Add the new user to the database using the repository.
// 	createdUser, err := us.UserRepo.AddUser(user)
// 	if err != nil {
// 		log.Printf("[ERROR - userService] Failed to add user to the database: %v", err)
// 		return "", err
// 	}
// 	return createdUser.ID.String(), nil
// }

func (us *UserService) AddUser(args dto.UserSignUp) (string, error) {
	log.Printf("[DEBUG] Inserting user with mobile: %s", args.Mobile)

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

	user := domain.User{
		ID:        uuid.New(),
		Mobile:    args.Mobile,
		Password:  hashedPassword,
		ExpiresAt: &expiry,
	}

	log.Printf("[DEBUG] User entity: %+v", user)

	createdUser, err := us.UserRepo.AddUser(user)
	if err != nil {
		log.Printf("[ERROR - UserService] Failed to add user to the database: %v", err)
		return "", err
	}

	return createdUser.ID.String(), nil
}

// FindUserByID retrieves a user by their unique ID.
// Returns the user or an error if the user could not be found or if an error occurs.
func (us *UserService) FindUserByID(userID uuid.UUID) (*domain.User, error) {
	user, err := us.UserRepo.FindUserByID(userID)
	if err != nil {
		log.Printf("[ERROR] Error finding user by ID %s: %v", userID, err)
		return nil, fmt.Errorf("could not find user by ID: %w", err)
	}
	return &user, nil
}

// FindUserByMobile retrieves a user by their mobile number.
// Returns the user or an error if the user could not be found or if an error occurs.
func (us *UserService) FindUserByMobile(mobile string) (*domain.User, error) {
	user, err := us.UserRepo.FindUserByMobile(mobile)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, nil // No user found, not an error
		}
		log.Printf("[ERROR] Error finding user by mobile: %v", err)
		return nil, err // Actual error occurred
	}
	return &user, nil
}

// FindUserByEmail retrieves a user by their email address.
// Returns the user or an error if the user could not be found or if an error occurs.
func (us *UserService) FindUserByEmail(email string) (*domain.User, error) {
	user, err := us.UserRepo.FindUserByEmail(email)
	if err != nil {
		log.Printf("[ERROR] Error finding user by email %s: %v", email, err)
		return nil, fmt.Errorf("could not find user by email: %w", err)
	}
	return &user, nil
}

// UpdateUser updates the details of an existing user based on the provided arguments.
// Returns the updated user or an error if the operation fails.
func (us *UserService) UpdateUser(userID uuid.UUID, arguments dto.UpdateUser) (*domain.User, error) {
	// Check if all fields in the arguments are nil
	if arguments.FirstName == nil &&
		arguments.LastName == nil &&
		arguments.Mobile == nil &&
		arguments.Email == nil &&
		arguments.Password == nil {
		log.Printf("[INFO] No arguments provided for user update")
		return nil, fmt.Errorf("nothing to update")
	}
	// Fetch the existing user from the database.
	user, err := us.UserRepo.FindUserByID(userID)
	if err != nil {
		log.Printf("[ERROR] Error finding user by ID %s: %v", userID, err)
		return nil, fmt.Errorf("could not find user by ID: %w", err)
	}

	// Update fields only if they are provided in the arguments.
	if arguments.FirstName != nil {
		user.FirstName = *arguments.FirstName
	}
	if arguments.LastName != nil {
		user.LastName = *arguments.LastName
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
			return nil, fmt.Errorf("failed to hash the password: %w", err)
		}
		user.Password = hashedPassword
	}
	// if arguments.ExpiresAt != nil {
	// 	user.ExpiresAt = arguments.ExpiresAt
	// }

	// Save the updated user to the database.
	updatedUser, err := us.UserRepo.UpdateUser(user)
	if err != nil {
		log.Printf("[ERROR] Failed to update user in the database: %v", err)
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	log.Printf("[INFO] User with ID %s successfully updated", userID)
	return &updatedUser, nil
}

// CountUsers counts the total number of users in the database.
// Returns the count or an error if the operation fails.
func (us *UserService) CountUsers() (int, error) {
	count, err := us.UserRepo.CountUsers()
	if err != nil {
		log.Printf("[ERROR] Error counting users: %v", err)
		return 0, fmt.Errorf("could not count users: %w", err)
	}
	return count, nil
}

// GetAllUsers retrieves all user records from the database.
// Returns a slice of users or an error if the operation fails.
func (us *UserService) GetAllUsers() ([]domain.User, error) {
	users, err := us.UserRepo.GetAllUsers()
	if err != nil {
		log.Printf("[ERROR] Error getting all users: %v", err)
		return nil, fmt.Errorf("could not get all users: %w", err)
	}
	return users, nil
}
