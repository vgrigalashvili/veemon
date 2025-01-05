// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID           uuid.UUID        `json:"id"`
	CreatedAt    pgtype.Timestamp `json:"created_at"`
	UpdatedAt    pgtype.Timestamp `json:"updated_at"`
	DeletedAt    pgtype.Timestamp `json:"deleted_at"`
	Role         string           `json:"role"`
	FirstName    pgtype.Text      `json:"first_name"`
	LastName     pgtype.Text      `json:"last_name"`
	Email        pgtype.Text      `json:"email"`
	Mobile       string           `json:"mobile"`
	PasswordHash string           `json:"password_hash"`
	Pin          pgtype.Int4      `json:"pin"`
	Verified     bool             `json:"verified"`
	UserType     string           `json:"user_type"`
	ExpiresAt    pgtype.Timestamp `json:"expires_at"`
}
