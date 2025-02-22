package dto

type AuthSignUp struct {
	Email  string `json:"email" validate:"required,email"`
	Secure bool   `json:"secure"`
}
type AuthSignIn struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}
