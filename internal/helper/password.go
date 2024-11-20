// Package helper provides utility functions for password generation, hashing, and comparison.
package helper

import (
	"crypto/rand"
	"log"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

// Constants for password generation.
const (
	// Characters used for password generation.
	lowercase = "abcdefghijklmnopqrstuvwxyz" // Lowercase alphabet characters.
	uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" // Uppercase alphabet characters.
	digits    = "0123456789"                 // Numerical digits.
	special   = "!@#$%^&*()-_=+<>?"          // Special characters for added complexity.

	// Combined set of all character types for general password generation.
	allChars = lowercase + uppercase + digits + special

	// Default length for generated passwords.
	passwordLength = 6 // Minimum length for passwords generated by the function.
)

// GeneratePassword generates a random password that includes lowercase, uppercase, digits, and special characters.
// Returns the generated password or an error if any random operation fails.
func GeneratePassword() (string, error) {
	password := make([]byte, passwordLength)

	charSets := []string{lowercase, uppercase, digits, special}
	for i := 0; i < len(charSets); i++ {
		char, err := randomCharFromSet(charSets[i])
		if err != nil {
			log.Printf("[ERROR] error generating character from set: %s, error: %v", charSets[i], err)
			return "", err
		}
		password[i] = char
	}

	for i := len(charSets); i < passwordLength; i++ {
		char, err := randomCharFromSet(allChars)
		if err != nil {
			log.Printf("[ERROR] error generating random character from allChars, error: %v", err)
			return "", err
		}
		password[i] = char
	}

	shuffle(password)

	return string(password), nil
}

// randomCharFromSet selects a random character from a given set and returns it.
// Returns an error if there is an issue with generating a random index.
func randomCharFromSet(set string) (byte, error) {
	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(set))))
	if err != nil {
		log.Printf("[ERROR] Error selecting random index from set: %s, error: %v", set, err)
		return 0, err
	}
	return set[index.Int64()], nil
}

// shuffle shuffles the password characters in place.
func shuffle(password []byte) {
	for i := len(password) - 1; i > 0; i-- {
		j, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		password[i], password[j.Int64()] = password[j.Int64()], password[i]
	}
}

// HashPassword hashes a plain text password using bcrypt.
// Returns the hashed password or an error if the hashing fails.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[ERROR] Error hashing password: %v", err)
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword compares a hashed password with a plain text password.
// Returns an error if the passwords do not match or if there is a comparison error.
func CheckPassword(hashedPassword, plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		log.Printf("[ERROR] Password comparison failed: %v", err)
		return err
	}
	return nil
}
