// Package token provides the definition and functionality for token payloads,
// including validation and creation.
package token

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

// Different types of errors returned by the VerifyToken function.
var (
	ErrInvalidToken = errors.New("token is invalid")  // Returned when a token is malformed or tampered with.
	ErrExpiredToken = errors.New("token has expired") // Returned when a token has exceeded its expiration time.
)

// Payload contains the payload data of the token, including metadata such as expiration and issuance times.
type Payload struct {
	UserID  uuid.UUID `json:"user_id"`      // Unique identifier
	TokenID uuid.UUID `json:"token_id"`     // Unique identifier for the token.
	Token   string    `json:"access_token"` // Access token associated with the token.
	// Email     string    `json:"email"`      // Email associated with the token.
	Role      string    `json:"role"`       // Role associated with the token (e.g., admin, user).
	IssuedAt  time.Time `json:"issued_at"`  // Time when the token was issued.
	ExpiredAt time.Time `json:"expired_at"` // Time when the token will expire.
}

// NewPayload creates a new token payload with a specific email, role, and duration.
// Returns a pointer to the Payload and an error if the operation fails.
func NewPayload(userID uuid.UUID, email, role string, duration time.Duration) (*Payload, error) {
	// Generate a unique ID for the token.
	tokenID, err := uuid.NewRandom()
	if err != nil {
		log.Printf("[ERROR] Failed to generate token ID: %v", err)
		return nil, err
	}

	// Create the token payload.
	payload := &Payload{
		UserID:  userID,
		TokenID: tokenID,
		Token:   "",
		// Email:     email,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

// Valid checks if the token payload is valid or not.
// Returns an error if the token is expired.
func (payload *Payload) Valid() error {
	// Check if the token has expired.
	if time.Now().After(payload.ExpiredAt) {
		log.Printf("[ERROR] Token has expired: ID=%s, ExpiredAt=%s", payload.TokenID, payload.ExpiredAt)
		return ErrExpiredToken
	}
	return nil
}
