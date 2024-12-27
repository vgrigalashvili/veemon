package repository

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx"
	"github.com/lib/pq"
	"github.com/vgrigalashvili/veemon/internal/domain"
	"gorm.io/gorm"
)

const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
)

var (
	ErrRecordNotFound = pgx.ErrNoRows

	ErrUniqueViolation = &pgconn.PgError{
		Code: UniqueViolation,
	}
)

type UserRepository interface {
	AddUser(user domain.User) (domain.User, error)
	FindUserByID(id uuid.UUID) (domain.User, error)
	FindUserByMobile(mobile string) (domain.User, error)
	FindUserByEmail(email string) (domain.User, error)
	UpdateUser(user domain.User) (domain.User, error)
	GetAllUsers() ([]domain.User, error)
	CountUsers() (int, error)
}

type userRepository struct {
	db *gorm.DB // Database connection instance.
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
func (ur *userRepository) AddUser(user domain.User) (domain.User, error) {
	log.Printf("[DEBUG] Attempting to insert user with mobile: %s", user.Mobile)

	if err := ur.db.Create(&user).Error; err != nil { // Pass by reference
		// Check for PostgreSQL duplicate key violation
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.Constraint == "uni_users_mobile" {
			log.Printf("[ERROR - UserRepository] User with mobile %s already exists", user.Mobile)
			return domain.User{}, fmt.Errorf("user with mobile %s already exists", user.Mobile)
		}

		log.Printf("[ERROR - UserRepository] Database error: %v", err)
		return domain.User{}, fmt.Errorf("database error: %w", err)
	}

	log.Printf("[DEBUG] User inserted successfully with mobile: %s", user.Mobile)
	return user, nil
}

func (ur *userRepository) FindUserByMobile(mobile string) (domain.User, error) {
	var user domain.User
	result := ur.db.Where("mobile = ?", mobile).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Printf("[DEBUG] User with mobile %s not found", mobile)
		return domain.User{}, ErrRecordNotFound
	} else if result.Error != nil {
		log.Printf("[ERROR] Error finding user by mobile %s: %v", mobile, result.Error)
		return domain.User{}, result.Error
	}
	return user, nil
}

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

func (ur *userRepository) UpdateUser(user domain.User) (domain.User, error) {
	if err := ur.db.Save(&user).Error; err != nil {
		log.Printf("[ERROR] Error updating user with ID %s: %v", user.ID, err)
		return domain.User{}, errors.New("unable to update user")
	}
	return user, nil
}

func (ur *userRepository) GetAllUsers() ([]domain.User, error) {
	var users []domain.User
	if err := ur.db.Find(&users).Error; err != nil {
		log.Printf("[ERROR] Error fetching users: %v", err)
		return nil, err
	}
	return users, nil
}

func (ur *userRepository) CountUsers() (int, error) {
	var count int64
	if err := ur.db.Model(&domain.User{}).Count(&count).Error; err != nil {
		log.Printf("[ERROR] Error while counting users: %v", err)
		return 0, err
	}
	return int(count), nil
}
