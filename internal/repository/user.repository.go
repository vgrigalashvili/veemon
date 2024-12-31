package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

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
	queries *db.Queries
}

// NewUserRepository creates a new instance of userRepository.
func NewUserRepository(db *db.Queries) UserRepository {
	return &userRepository{
		queries: db,
	}
}

// CreateUser inserts a new user into the database.
func (ur *userRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	dbUser, err := ur.queries.CreateUser(ctx, db.CreateUserParams{
		ID:           user.ID,
		Role:         user.Role,
		FirstName:    pgtype.Text{String: user.FirstName},
		LastName:     pgtype.Text{String: user.LastName},
		Email:        pgtype.Text{String: user.Email},
		Mobile:       user.Mobile,
		PasswordHash: user.Password,
		Code:         pgtype.Int4{Int32: int32(user.Code)},
		Verified:     user.Verified,
		UserType:     user.UserType,
		ExpiresAt:    pgtype.Timestamp{Time: *user.ExpiresAt},
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
	dbUsers, err := ur.queries.GetAllUsersPaginated(ctx, sqlc.GetAllUsersPaginatedParams{
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
		Role:      dbUser.Role,                      // Convert pgtype.Text to string
		FirstName: dbUser.FirstName,                 // Convert pgtype.Text to string
		LastName:  dbUser.LastName,                  // Convert pgtype.Text to string
		Email:     dbUser.Email,                     // Convert pgtype.Text to string
		Mobile:    dbUser.Mobile,                    // Convert pgtype.Text to string
		Password:  dbUser.Password,                  // Convert pgtype.Text to string
		Code:      dbUser.Code,                      // Convert pgtype.Int4 to int
		Verified:  dbUser.Verified,                  // Convert pgtype.Bool to bool
		UserType:  dbUser.UserType,                  // Convert pgtype.Text to string
		ExpiresAt: extractTimePtr(dbUser.ExpiresAt), // Convert pgtype.Timestamp to *time.Time
	}
}
