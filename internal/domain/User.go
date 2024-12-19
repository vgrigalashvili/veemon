// Package domain defines the structure for the User entity used in the application.
// This struct is designed for use with GORM, an ORM library for Go.
package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system with relevant fields and metadata for GORM.
type User struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time      // The timestamp when the user record was created.
	UpdatedAt time.Time      // The timestamp when the user record was last updated.
	DeletedAt gorm.DeletedAt `gorm:"index"` // Soft delete field with an index for querying.

	Role      string     `json:"role" gorm:"type:varchar(255);default:'user'"`                          // The role assigned to the user (e.g., admin, user).
	FirstName string     `json:"first_name"`                                                            // The first name of the user.
	LastName  string     `json:"last_name"`                                                             // The last name of the user.
	Email     string     `json:"email" gorm:"index;unique;default:NULL"`                                // The email of the user, unique.
	Mobile    string     `json:"mobile" gorm:"index;unique;not null" validate:"required,len=9,numeric"` // The mobile phone number of the user, unique and non-null.
	Password  string     `json:"password"`                                                              // The password for the user, should be hashed.
	Code      int        `json:"code"`                                                                  // An optional verification or identification code.
	Verified  bool       `json:"verified" gorm:"default:false"`                                         // Indicates if the user's email or account is verified.
	UserType  string     `json:"user_type" gorm:"default:trial"`                                        // Indicates the type of user, default is "trial".
	ExpiresAt *time.Time `json:"expires_at" gorm:"default:NULL"`                                        // Optional expiration date for user-specific tokens or sessions.
}
