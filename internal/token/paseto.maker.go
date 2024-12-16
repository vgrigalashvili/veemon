// Package token provides functionality for creating and verifying tokens using PASETO (Platform-Agnostic Security Tokens).
package token

import (
	"fmt"
	"log"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

// PasetoMaker is a struct that implements the Maker interface using PASETO for token generation and verification.
type PasetoMaker struct {
	paseto       *paseto.V2 // Paseto V2 instance for handling tokens.
	symmetricKey []byte     // Symmetric key used for encrypting and decrypting tokens.
}

// NewPasetoMaker creates a new PasetoMaker with the given symmetric key.
// Returns an error if the key size is invalid.
func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		log.Printf("[ERROR] Invalid key size: expected %d characters, got %d", chacha20poly1305.KeySize, len(symmetricKey))
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

// CreateToken generates a new token for a given email and role with a specific duration.
// Returns the token string, payload, and any error encountered during token creation.
func (maker *PasetoMaker) CreateToken(userID uuid.UUID, email, role string, duration time.Duration) (*Payload, error) {
	// Create the payload with email, role, and expiration details.
	payload, err := NewPayload(userID, email, role, duration)
	if err != nil {
		log.Printf("[ERROR] Failed to create payload: %v", err)
		return payload, err
	}
	// Encrypt the payload into a PASETO token.
	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	if err != nil {
		log.Printf("[ERROR] Failed to encrypt token: %v", err)
		return payload, err
	}

	payload.Token = token

	return payload, err
}

// VerifyToken decrypts and validates a given token string.
// Returns the payload if valid, or an error if the token is invalid or expired.
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	// Decrypt the token and populate the payload.
	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		log.Printf("[ERROR] Failed to decrypt token: %v", err)
		return nil, ErrInvalidToken
	}

	// Validate the payload for expiration and other checks.
	err = payload.Valid()
	if err != nil {
		log.Printf("[ERROR] Token validation failed: %v", err)
		return nil, err
	}

	return payload, nil
}
