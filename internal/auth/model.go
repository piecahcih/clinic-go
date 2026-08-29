package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"            db:"id"`
	Email        string    `json:"email"         db:"email"`
	PasswordHash string    `json:"-"             db:"password_hash"`
	FirstName    string    `json:"firstName"     db:"first_name"`
	LastName     string    `json:"lastName"      db:"last_name"`
	BirthDate    time.Time `json:"birthDate"     db:"birth_date"`
	Gender       string    `json:"gender"        db:"gender"`
	Role         string    `json:"role"          db:"role"`
	CreatedAt    time.Time `json:"createdAt"     db:"created_at"`
}

var (
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
)
