package db

import "errors"

var (
	ErrNotFound     = errors.New("record not found")
	ErrUnique       = errors.New("unique constraint violation")
	ErrConflict     = errors.New("conflicting key value violates unique constraint")
	ErrNoRows       = errors.New("no rows affected")
	ErrInvalid      = errors.New("invalid input syntax for type")
	ErrInternal     = errors.New("internal server error")
	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidPhone = errors.New("invalid phone number")
	ErrInvalidRole  = errors.New("invalid role")
	ErrInvalidPin   = errors.New("invalid pin")
	ErrInvalidType  = errors.New("invalid user type")
	ErrInvalidUUID  = errors.New("invalid UUID")
)
