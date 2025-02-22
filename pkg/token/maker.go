package token

import (
	"time"

	"github.com/google/uuid"
)

// Maker is an interface for managing tokens
type Maker interface {
	// CreateToken creates a new token for a specific username and duration
	CreateToken(userID uuid.UUID, email, role string, duration time.Duration) (*Payload, error)

	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
