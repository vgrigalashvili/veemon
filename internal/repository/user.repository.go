package repository

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"

	"github.com/vgrigalashvili/veemon/internal/domain"
	db "github.com/vgrigalashvili/veemon/internal/repository/sqlc"
)

var (
	ErrNoRows             = errors.New("no rows found")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user with this mobile already exists")
	ErrEmailAlreadyExists = errors.New("user with this email already exists")
	ErrPasswordMismatch   = errors.New("password mismatch")
	ErrUserNotVerified    = errors.New("user not verified")
	ErrUserExpired        = errors.New("user expired")
	ErrUserDeleted        = errors.New("user deleted")
	ErrUserNotDeleted     = errors.New("user not deleted")
	ErrUserNotUpdated     = errors.New("user not updated")
	ErrUserNotCreated     = errors.New("user not created")
)

// UserRepository defines the methods for interacting with the user data store.
type (
	UserRepository interface {
		Create(ctx context.Context, user domain.User) (*domain.User, error)
		Read(ctx context.Context, id uuid.UUID) (*domain.User, error)
		// Update(ctx context.Context, user domain.User) (*domain.User, error)
		// Delete(ctx context.Context, id uuid.UUID) error
	}
)

// userRepository implements the UserRepository interface using sqlc.
type userRepository struct {
	queries *db.Queries
}

// NewUserRepository creates a new instance of userRepository.
// It requires a *db.Queries instance to be passed in.
func NewUserRepository(q *db.Queries) UserRepository {
	if q == nil {
		log.Fatalf("[FATAL] queries cannot be nil")
	}
	return &userRepository{queries: q}
}

func (ur *userRepository) Create(ctx context.Context, user domain.User) (*domain.User, error) {
	// Use the SQLC generated method instead of recursive call
	dbUser, err := ur.queries.CreateUser(ctx, domainToDBUser(user))
	if err != nil {
		return nil, err
	}
	return dbToDomainUser(dbUser), nil
}

func (ur *userRepository) Read(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := ur.queries.Read(ctx, id)
	if err != nil {

		return nil, err
	}
	return &domain.User{
		ID:             user.ID,
		FirstName:      *user.FirstName,
		LastName:       *user.LastName,
		Email:          user.Email,
		Email_verified: user.EmailVerified,
	}, nil
}
func domainToDBUser(u domain.User) db.CreateUserParams {
	return db.CreateUserParams{
		ID:        u.ID,
		FirstName: &u.FirstName,
		LastName:  &u.LastName,
		Email:     u.Email,
		Password:  u.Password,
	}
}

func dbToDomainUser(u db.User) *domain.User {
	return &domain.User{
		ID:        u.ID,
		FirstName: *u.FirstName,
		LastName:  *u.LastName,
		Email:     u.Email,
		Password:  u.Password,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt.Time,
	}
}
