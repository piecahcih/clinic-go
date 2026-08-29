package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo   Repository
	secret string
	ttl    time.Duration
}

func NewService(r Repository, secret string, ttl time.Duration) *Service {
	return &Service{repo: r, secret: secret, ttl: ttl}
}

func (s *Service) Register(ctx context.Context, u User, password string) (*User, error) {
	exists, err := s.repo.EmailExists(ctx, u.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	u.ID = uuid.New()
	u.PasswordHash = string(hash)
	if u.Role == "" {
		u.Role = "patient"
	}

	if err := s.repo.CreateUser(ctx, &u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	u, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials // don't leak "no such user"
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials // same error, deliberately
	}

	claims := jwt.MapClaims{
		"sub":  u.ID.String(),
		"role": u.Role,
		"exp":  time.Now().Add(s.ttl).Unix(),
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}
