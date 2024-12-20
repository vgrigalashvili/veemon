package dto

// UserSignUp represents the data required for a user to sign up.
type AuthSignUp struct {
	Mobile string `json:"mobile" validate:"required,min=9"` // Mobile number, required and with a minimum length of 9 characters.
}

// UserSignIn represents the data required for a user to sign in.
type AuthSignIn struct {
	Mobile   string `json:"mobile" validate:"required,min=9"`   // Mobile number, required and with a minimum length of 9 characters.
	Password string `json:"password" validate:"required,min=6"` // Password, required and with a minimum length of 6 characters.
}
