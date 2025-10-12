package services

import (
	"crypto/rand"
	"crypto/subtle"
	"fmt"
	"golang.org/x/crypto/argon2"
)

func HashPass(password string) ([]byte, error) {
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("hashing error: %w", err)
	}
	hashedPass := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return append(salt, hashedPass...), nil
}

func VerifyPassword(password string, hashedPassword []byte) (bool, error) {
	if len(hashedPassword) < 32 {
		return false, fmt.Errorf("invalid hashed password length: %v", len(hashedPassword))
	}
	salt := hashedPassword[:32]
	hash := hashedPassword[32:]

	newHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	return subtle.ConstantTimeCompare(hash, newHash) == 1, nil
}
