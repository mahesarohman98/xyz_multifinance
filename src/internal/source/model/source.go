package model

import (
	"golang.org/x/crypto/bcrypt"
)

type Source struct {
	ID         string `db:"source_id"`
	SecretHash string `db:"secret_hash"`
	Category   string `db:"category"`
	Name       string `db:"name"`
	Email      string `db:"email"`
}

func (s *Source) GenerateHashFromSecret(secret string) error {
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	s.SecretHash = string(hashedSecret)
	return nil
}

func (s Source) Validate(secret string) error {
	return bcrypt.CompareHashAndPassword([]byte(s.SecretHash), []byte(secret))
}
