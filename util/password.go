package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// returns the bcrypt hash of a password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %s", err)
	}

	return string(hash), nil
}

// checks if the provided password is correct or not
func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}