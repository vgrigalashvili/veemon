// Package repository provides data access to user entities using GORM and PostgreSQL.
package repository

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx"
	"github.com/vgrigalashvili/veemon/internal/domain"
	"gorm.io/gorm"
)

const (
	// PostgreSQL error codes for constraint violations.
	ForeignKeyViolation = "23503" // Code for foreign key violation errors.
	UniqueViolation     = "23505" // Code for unique constraint violation errors.
)

var (
	// ErrRecordNotFound is returned when a database query does not find a record.
	ErrRecordNotFound = pgx.ErrNoRows

	// ErrUniqueViolation represents a unique constraint violation error.
	ErrUniqueViolation = &pgconn.PgError{
		Code: UniqueViolation,
	}
)

// UserRepository defines the interface for user-related database operations.
// It allows for consistent method signatures for implementing user operations.
type UserRepository interface {
	AddUser(user domain.User) (domain.User, error)       // Adds a new user to the database.
	FindUserByMobile(mobile string) (domain.User, error) // Finds a user by their mobile number.
	FindUserByEmail(email string) (domain.User, error)   // Finds a user by their email address.
	FindUserByID(id uuid.UUID) (domain.User, error)      // Finds a user by their ID.
	UpdateUser(user domain.User) (domain.User, error)    // Updates an existing user's details.
	GetAllUsers() ([]domain.User, error)                 // Retrieves all users from the database.
	CountUsers() (int, error)                            // Counts the total number of users.
}

// userRepository is the concrete implementation of UserRepository.
type userRepository struct {
	db *gorm.DB // Database connection instance.
}

// NewUserRepository creates and returns a new userRepository instance.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// AddUser inserts a new user record into the database.
// Returns the created user or an error if the operation fails.
func (ur *userRepository) AddUser(user domain.User) (domain.User, error) {
	log.Printf("[DEBUG] Adding user with mobile: %s", user.Mobile)
	if err := ur.db.Create(&user).Error; err != nil {
		log.Printf("[ERROR] Unable to create new user: %v", err)
		// Check if the error is a unique constraint violation.
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == UniqueViolation {
			log.Printf("[ERROR] Unique violation for user with mobile: %s", user.Mobile)
			return domain.User{}, errors.New("user with this mobile already exists")
		}
		return domain.User{}, errors.New("unable to create new user")
	}
	log.Printf("[DEBUG] User created successfully with ID: %s", user.ID)
	return user, nil
}

// FindUserByMobile retrieves a user by their mobile number.
// Returns the user or an error if not found or if another error occurs.
func (ur *userRepository) FindUserByMobile(mobile string) (domain.User, error) {
	var user domain.User
	result := ur.db.Where("mobile = ?", mobile).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Printf("[ERROR] User with mobile %s not found", mobile)
		return domain.User{}, ErrRecordNotFound
	} else if result.Error != nil {
		log.Printf("[ERROR] Error finding user by mobile %s: %v", mobile, result.Error)
		return domain.User{}, result.Error
	}

	return user, nil
}

// FindUserByEmail retrieves a user by their email address.
// Returns the user or an error if not found or if another error occurs.
func (ur *userRepository) FindUserByEmail(email string) (domain.User, error) {
	var user domain.User
	result := ur.db.Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Printf("[ERROR] User with email %s not found", email)
		return domain.User{}, ErrRecordNotFound
	} else if result.Error != nil {
		log.Printf("[ERROR] Error finding user by email %s: %v", email, result.Error)
		return domain.User{}, result.Error
	}

	return user, nil
}

// FindUserByID retrieves a user by their ID.
// Returns the user or an error if not found or if another error occurs.
func (ur *userRepository) FindUserByID(id uuid.UUID) (domain.User, error) {
	var user domain.User
	result := ur.db.Where("id = ?", id).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Printf("[ERROR] User with ID %s not found", id)
		return domain.User{}, ErrRecordNotFound
	} else if result.Error != nil {
		log.Printf("[ERROR] Error finding user by ID %s: %v", id, result.Error)
		return domain.User{}, result.Error
	}
	return user, nil
}

// UpdateUser updates an existing user's details in the database.
// Returns the updated user or an error if the operation fails.
func (ur *userRepository) UpdateUser(user domain.User) (domain.User, error) {
	if err := ur.db.Save(&user).Error; err != nil {
		log.Printf("[ERROR] Error updating user with ID %s: %v", user.ID, err)
		return domain.User{}, errors.New("unable to update user")
	}
	return user, nil
}

// GetAllUsers retrieves all user records from the database.
// Returns a slice of users or an error if the operation fails.
func (ur *userRepository) GetAllUsers() ([]domain.User, error) {
	var users []domain.User
	if err := ur.db.Find(&users).Error; err != nil {
		log.Printf("[ERROR] Error fetching users: %v", err)
		return nil, err
	}
	return users, nil
}

// CountUsers counts the total number of users in the database.
// Returns the count or an error if the operation fails.
func (ur *userRepository) CountUsers() (int, error) {
	var count int64
	if err := ur.db.Model(&domain.User{}).Count(&count).Error; err != nil {
		log.Printf("[ERROR] Error while counting users: %v", err)
		return 0, err
	}
	return int(count), nil
}
