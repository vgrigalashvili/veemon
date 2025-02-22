package validator

import (
	"github.com/go-playground/validator/v10"
)

// CustomValidator struct
type CustomValidator struct {
	Validator *validator.Validate
}

// NewValidator initializes validator with custom rules
func NewValidator() *CustomValidator {
	v := validator.New()
	// Add custom validation rules if needed
	return &CustomValidator{Validator: v}
}

// ValidateStruct validates a struct and returns formatted errors
func (cv *CustomValidator) ValidateStruct(s interface{}) error {
	return cv.Validator.Struct(s)
}
