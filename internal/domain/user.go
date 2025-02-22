package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time // The timestamp when the user record was created.
	UpdatedAt time.Time // The timestamp when the user record was last updated.
	DeletedAt time.Time // Soft delete field with an index for querying.

	FirstName      string `json:"first_name"`     // The first name of the user.
	LastName       string `json:"last_name"`      // The last name of the user.
	Email          string `json:"email"`          // The email of the user, unique.
	Password       string `json:"password"`       // The password for the user, should be hashed.
	Role           string `json:"role"`           // The role assigned to the user (e.g., admin, user).
	Email_verified bool   `json:"email_verified"` // Indicates if the user's email or account is verified.
}
