// Package dto contains Data Transfer Objects for handling user input and output,
// such as user sign-in and sign-up requests.
package dto

// UserSignIn represents the data required for a user to sign in.
type UserSignIn struct {
	Mobile   string `json:"mobile" validate:"required,min=9"`   // Mobile number, required and with a minimum length of 9 characters.
	Password string `json:"password" validate:"required,min=6"` // Password, required and with a minimum length of 6 characters.
}

// UserSignUp represents the data required for a user to sign up.
type UserSignUp struct {
	Mobile string `json:"mobile" validate:"required,min=9"` // Mobile number, required and with a minimum length of 9 characters.
}

type CreateUser struct {
	FirstName string `json:"first_name" validate:"required,min=2"` // First name, required.
	LastName  string `json:"last_name" validate:"required,min=3"`  // Last name, required.
	Mobile    string `json:"mobile" validate:"required,min=9"`     // Mobile number, required and with a minimum length of 9 characters.
	Email     string `json:"email" validate:"required,email"`      // Email address, required.
	Password  string `json:"password" validate:"omitempty,min=6"`  // Password, required and with a minimum length of 6 characters.
	Role      string `json:"role" validate:"required"`
}
type UpdateUser struct {
	FirstName *string `json:"first_name" validate:"omitempty"`     // First name (optional)
	LastName  *string `json:"last_name" validate:"omitempty"`      // Last name (optional)
	Mobile    *string `json:"mobile" validate:"omitempty,min=9"`   // Mobile number (optional, min length 9 if provided)
	Email     *string `json:"email" validate:"omitempty,email"`    // Email address (optional, must be a valid email if provided)
	Password  *string `json:"password" validate:"omitempty,min=6"` // Password (optional, min length 6 if provided)
	Role      *string `json:"role" validate:"omitempty"`

	// ExpiresAt *time.Time `json:"expires_at" validate:"omitempty"`           // Expiration date (optional)
}
