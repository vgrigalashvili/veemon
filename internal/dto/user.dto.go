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
