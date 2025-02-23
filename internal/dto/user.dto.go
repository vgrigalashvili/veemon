package dto

type CreateUser struct {
	FirstName string `json:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" validate:"required,min=3"`
	Mobile    string `json:"mobile" validate:"required,min=9"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"omitempty,min=6"`
	Role      string `json:"role" validate:"required"`
}
type UpdateUser struct {
	FirstName *string `json:"first_name" validate:"omitempty"`
	LastName  *string `json:"last_name" validate:"omitempty"`
	Type      *string `json:"type" validate:"omitempty"`
	Role      *string `json:"role" validate:"omitempty"`
	Mobile    *string `json:"mobile" validate:"omitempty,min=9"`
	Email     *string `json:"email" validate:"omitempty,email"`
	Password  *string `json:"password" validate:"omitempty,min=6"`

	// ExpiresAt *time.Time `json:"expires_at" validate:"omitempty"`
}
