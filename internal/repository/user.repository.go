package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	db "github.com/vgrigalashvili/veemon/internal/db/sqlc"
	"github.com/vgrigalashvili/veemon/internal/domain"
)

// UserRepository defines the methods for interacting with the user data store.
type UserRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetAllUsers(ctx context.Context, limit, offset int) ([]domain.User, error)
	CheckUserExistsByMobile(ctx context.Context, mobile string) (bool, error)
	CheckUserExistsByEmail(ctx context.Context, email string) (bool, error)
	UpdateUser(ctx context.Context, user domain.User) (domain.User, error)
	SoftDeleteUser(ctx context.Context, id uuid.UUID) error
	HardDeleteUser(ctx context.Context, id uuid.UUID) error
}

// userRepository implements the UserRepository interface using sqlc.
type userRepository struct {
	connPool *pgx.ConnPool
	queries  *db.Queries
}

// NewUserRepository creates a new instance of userRepository.
func NewUserRepository(connPool *pgx.ConnPool, queries *db.Queries) UserRepository {
	return &userRepository{
		connPool: connPool,
		queries:  queries,
	}
}

// CreateUser inserts a new user into the database.
func (ur *userRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	err := ur.queries.CreateUser(ctx, db.CreateUserParams{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Role:      user.Role,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Mobile:    user.Mobile,
		Password:  user.Password,
		Code:      user.Code,
		Verified:  user.Verified,
		UserType:  user.UserType,
		ExpiresAt: user.ExpiresAt,
	})
	if err != nil {
		return domain.User{}, err
	}
	return mapDBUserToDomainUser(dbUser), nil
}

// GetUserByID retrieves a user by ID, excluding soft-deleted users.
func (ur *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	dbUser, err := ur.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}
	return mapDBUserToDomainUser(dbUser), nil
}

// GetAllUsers retrieves all users with pagination, excluding soft-deleted users.
func (ur *userRepository) GetAllUsers(ctx context.Context, limit, offset int) ([]domain.User, error) {
	dbUsers, err := ur.queries.GetAllUsersPaginated(ctx, db.GetAllUsersPaginatedParams{
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

// CheckUserExistsByMobile checks if a user exists with the given mobile number.
func (ur *userRepository) CheckUserExistsByMobile(ctx context.Context, mobile string) (bool, error) {
	_, err := ur.queries.CheckUserExistsByMobile(ctx, mobile)
	if errors.Is(err, db.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

// CheckUserExistsByEmail checks if a user exists with the given email address.
func (ur *userRepository) CheckUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	_, err := ur.queries.CheckUserExistsByEmail(ctx, email)
	if errors.Is(err, db.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

// UpdateUser updates an existing user in the database.
func (ur *userRepository) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {

	err := ur.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:        user.ID,
		UpdatedAt: user.UpdatedAt,
		Role:      user.Role,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Mobile:    user.Mobile,
		Password:  user.Password,
		Code:      user.Code,
		Verified:  user.Verified,
		UserType:  user.UserType,
		ExpiresAt: user.ExpiresAt,
	})
	if err != nil {
		return domain.User{}, err
	}
	return mapDBUserToDomainUser(dbUser), nil
}

// SoftDeleteUser soft deletes a user by setting the deleted_at timestamp.
func (ur *userRepository) SoftDeleteUser(ctx context.Context, id uuid.UUID) error {
	return ur.queries.SoftDeleteUser(ctx, id)
}

// HardDeleteUser permanently deletes a user from the database.
func (ur *userRepository) HardDeleteUser(ctx context.Context, id uuid.UUID) error {
	return ur.queries.HardDeleteUser(ctx, id)
}

func mapDBUserToDomainUser(dbUser db.User) domain.User {
	return domain.User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt.Time,            // Convert pgtype.Timestamp to time.Time
		UpdatedAt: dbUser.UpdatedAt.Time,            // Convert pgtype.Timestamp to time.Time
		Role:      toString(dbUser.Role),            // Convert pgtype.Text to string
		FirstName: toString(dbUser.FirstName),       // Convert pgtype.Text to string
		LastName:  toString(dbUser.LastName),        // Convert pgtype.Text to string
		Email:     toString(dbUser.Email),           // Convert pgtype.Text to string
		Mobile:    toString(dbUser.Mobile),          // Convert pgtype.Text to string
		Password:  toString(dbUser.Password),        // Convert pgtype.Text to string
		Code:      toInt(dbUser.Code),               // Convert pgtype.Int4 to int
		Verified:  toBool(dbUser.Verified),          // Convert pgtype.Bool to bool
		UserType:  toString(dbUser.UserType),        // Convert pgtype.Text to string
		ExpiresAt: extractTimePtr(dbUser.ExpiresAt), // Convert pgtype.Timestamp to *time.Time
	}
}

// extractTimePtr handles nullable pgtype.Timestamp.
func extractTimePtr(ts pgtype.Timestamp) *time.Time {
	if ts.Status == pgtype.Null {
		return nil
	}
	return &ts.Time
}

// toPgTimestamp converts time.Time to pgtype.Timestamp.
func toPgTimestamp(t time.Time) pgtype.Timestamp {
	var ts pgtype.Timestamp
	ts.Set(t)
	return ts
}

func toBool(b pgtype.Bool) bool {
	if b.Status == pgtype.Null {
		return false
	}
	return b.Bool
}

// toString converts pgtype.Text to string.
func toString(pgText pgtype.Text) string {
	if pgText.Status == pgtype.Null {
		return ""
	}
	return pgText.String
}

// toInt converts pgtype.Int4 to int.
func toInt(pgInt pgtype.Int4) int {
	if pgInt.Status == pgtype.Null {
		return 0
	}
	return int(pgInt.Int32)
}
