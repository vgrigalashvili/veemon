package handler

import "errors"

var (
	// user handler errors
	ErrInvalidUserIDFormat = errors.New("Invalid user ID format")
	ErrInvalidEmail        = errors.New("invalid email format")

	ErrEmailQueryParamRequired = errors.New("query parameter required: email")
	ErrUnverified              = errors.New("unverified user")
	ErrUnauthorized            = errors.New("unauthorized")
	ErrNotFound                = errors.New("not found")
	ErrInvalidOrExpiredToken   = errors.New("invalid or expired token")
	ErrInvalidMethod           = errors.New("invalid method")
	ErrInvalidQueryParam       = errors.New("invalid query parameter")
	ErrInvalidRequestJSON      = errors.New("invalid JSON body in request")
	ErrValidationField         = errors.New("validation field")
	ErrUniqueMobileComplaint   = errors.New("mobile already taken")
)
