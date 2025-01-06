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

	Role      string     `json:"role"`       // The role assigned to the user (e.g., admin, user).
	FirstName string     `json:"first_name"` // The first name of the user.
	LastName  string     `json:"last_name"`  // The last name of the user.
	Email     string     `json:"email"`      // The email of the user, unique.
	Mobile    string     `json:"mobile"`     // The mobile phone number of the user, unique and non-null.
	Password  string     `json:"password"`   // The password for the user, should be hashed.
	Pin       int        `json:"pin"`        // An optional verification or identification code.
	Verified  bool       `json:"verified"`   // Indicates if the user's email or account is verified.
	UserType  string     `json:"user_type"`  // Indicates the type of user, default is "trial".
	ExpiresAt *time.Time `json:"expires_at"` // Optional expiration date for user-specific tokens or sessions.
}
