package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/vgrigalashvili/veemon/internal/db/sqlc"
	"github.com/vgrigalashvili/veemon/internal/domain"
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
		UserModifiers
		UserGetters
	}
	UserModifiers interface {
		CreateUser(ctx context.Context, user domain.User) (domain.User, error)
		UpdateUser(ctx context.Context, user domain.User) (domain.User, error)
		SoftDeleteUser(ctx context.Context, id uuid.UUID) error
		UserExpiresAt(ctx context.Context, id uuid.UUID, expiresAt time.Time) error
		SetupUserRole(ctx context.Context, id uuid.UUID, role string) error
	}
	UserGetters interface {
		GetUserByMobile(ctx context.Context, mobile string) (domain.User, error)
		GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
		GetUserRole(ctx context.Context, id uuid.UUID) (string, error)
		GetAllUsers(ctx context.Context, limit, offset int) ([]domain.User, error)
		CheckUserExistsByMobile(ctx context.Context, mobile string) bool
		CheckUserExistsByEmail(ctx context.Context, email string) bool
	}
)

// userRepository implements the UserRepository interface using sqlc.
type userRepository struct {
	queries *db.Queries
}

// NewUserRepository creates a new instance of userRepository.
// It requires a *db.Queries instance to be passed in.
func NewUserRepository(q *db.Queries) *userRepository {
	if q == nil {
		log.Fatalf("[FATAL] queries cannot be nil")
	}
	return &userRepository{queries: q}
}

// CreateUser inserts a new user into the database.
// Returns the created user or an error if the operation failed.
func (ur *userRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	dbUser, err := ur.queries.AddUser(ctx, db.AddUserParams{
		ID:           user.ID,
		Role:         user.Role,
		FirstName:    pgtype.Text{String: user.FirstName, Valid: true},
		LastName:     pgtype.Text{String: user.LastName, Valid: true},
		Email:        pgtype.Text{String: user.Email, Valid: true},
		Mobile:       user.Mobile,
		PasswordHash: user.Password,
		Pin:          pgtype.Int4{Int32: int32(user.Pin), Valid: true},
		Verified:     user.Verified,
		UserType:     user.UserType,
		ExpiresAt:    pgtype.Timestamp{Time: *user.ExpiresAt, Valid: true},
	})
	if err != nil {
		return domain.User{}, err
	}
	return mapDBUserToDomainUser(dbUser), nil
}

// GetUserByID retrieves a user by ID, excluding soft-deleted users.
// Returns the user or an error if the user could not be found or if an error occurs.
func (ur *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	dbUser, err := ur.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return domain.User{}, ErrUserNotFound
		}
		return domain.User{}, err
	}
	return mapDBUserToDomainUser(dbUser), nil
}

// GetAllUsers retrieves all users with pagination, excluding soft-deleted users.
// Returns a slice of users or an error if the operation failed.
func (ur *userRepository) GetAllUsers(ctx context.Context, limit, offset int) ([]domain.User, error) {
	dbUsers, err := ur.queries.ListUsers(ctx, db.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}
	var users []domain.User
	for _, dbUser := range dbUsers {
		users = append(users, mapDBUserToDomainUser(dbUser))
	}
	return users, nil
}

// GetUserByMobile retrieves a user by mobile number.
// Returns the user or an error if the user could not be found or if an error occurs.
func (ur *userRepository) GetUserByMobile(ctx context.Context, mobile string) (domain.User, error) {
	dbUser, err := ur.queries.GetUserByMobile(ctx, mobile)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return domain.User{}, ErrUserNotFound
		}
		return domain.User{}, err
	}
	return mapDBUserToDomainUser(dbUser), nil
}

// CheckUserExistsByMobile checks if a user exists with the given mobile number.
// Returns true if the user exists, false otherwise.
func (ur *userRepository) CheckUserExistsByMobile(ctx context.Context, mobile string) bool {
	_, err := ur.queries.GetUserByMobile(ctx, mobile)
	if err != nil {
		// Check if the error is due to no rows being found
		if errors.Is(err, pgx.ErrNoRows) {
			return false
		}
		// Log unexpected errors
		log.Printf("[ERROR] something went wrong in *userRepository: CheckUserExistsByMobile %v", err)
	}
	return true
}

// CheckUserExistsByEmail checks if a user exists with the given email address.
// Returns true if the user exists, false otherwise.
func (ur *userRepository) CheckUserExistsByEmail(ctx context.Context, email string) bool {
	_, err := ur.queries.GetUserByEmail(ctx, pgtype.Text{String: email})
	if err != nil && errors.Is(err, db.ErrNoRows) {
		return false
	} else if err != nil {
		log.Printf("[ERROR] something went wrong in *userRepository: CheckUserExistsByEmail %v", err)
	}
	return true
}

// UpdateUser updates an existing user in the database.
// Returns the updated user or an error if the operation failed.
func (ur *userRepository) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	err := ur.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:           user.ID,
		Role:         user.Role,
		FirstName:    pgtype.Text{String: user.FirstName},
		LastName:     pgtype.Text{String: user.LastName},
		Email:        pgtype.Text{String: user.Email},
		Mobile:       user.Mobile,
		PasswordHash: user.Password,
		Pin:          pgtype.Int4{Int32: int32(user.Pin)},
		Verified:     user.Verified,
		UserType:     user.UserType,
		ExpiresAt:    pgtype.Timestamp{Time: *user.ExpiresAt},
	})
	if err != nil {
		return domain.User{}, err
	}
	return ur.GetUserByID(ctx, user.ID)
}

// SoftDeleteUser soft deletes a user by setting the deleted_at timestamp.
func (ur *userRepository) SoftDeleteUser(ctx context.Context, id uuid.UUID) error {
	return ur.queries.SoftDeleteUser(ctx, id)
}

// UserExpiresAt sets the expiration date for a user.
func (ur *userRepository) UserExpiresAt(ctx context.Context, id uuid.UUID, expiresAt time.Time) error {
	return ur.queries.UserExpiresAt(ctx, db.UserExpiresAtParams{
		ID:        id,
		ExpiresAt: pgtype.Timestamp{Time: expiresAt, Valid: true},
	})
}

// GetUserRole retrieves the role of a user by ID.
func (ur *userRepository) GetUserRole(ctx context.Context, id uuid.UUID) (string, error) {
	role, err := ur.queries.GetUserRole(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return "", ErrUserNotFound
		}
		return "", err
	}
	return role, nil
}

// SetupUserRole sets the role of a user.
func (ur *userRepository) SetupUserRole(ctx context.Context, id uuid.UUID, role string) error {
	return ur.queries.SetupUserRole(ctx, db.SetupUserRoleParams{
		ID:   id,
		Role: role,
	})
}

// mapDBUserToDomainUser maps a db.User to a domain.User.
func mapDBUserToDomainUser(dbUser db.User) domain.User {
	return domain.User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
		Role:      dbUser.Role,
		FirstName: dbUser.FirstName.String,
		LastName:  dbUser.LastName.String,
		Email:     dbUser.Email.String,
		Mobile:    dbUser.Mobile,
		Password:  dbUser.PasswordHash,
		Pin:       int(dbUser.Pin.Int32),
		Verified:  dbUser.Verified,
		UserType:  dbUser.UserType,
		ExpiresAt: extractTimePtr(dbUser.ExpiresAt),
	}
}

// extractTimePtr safely converts pgtype.Timestamp to *time.Time.
func extractTimePtr(ts pgtype.Timestamp) *time.Time {
	if ts.Valid {
		return &ts.Time
	}
	return nil
}
