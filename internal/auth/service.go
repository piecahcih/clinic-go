package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo       Repository
	secret     string
	ttl        time.Duration
	refreshTTL time.Duration
}

func NewService(r Repository, secret string, ttl, refreshTTL time.Duration) *Service {
	return &Service{repo: r, secret: secret, ttl: ttl, refreshTTL: refreshTTL}
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

func (s *Service) signJWT(u *User) (string, error) {
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

func (s *Service) issueTokens(ctx context.Context, u *User) (access, refresh string, err error) {
	access, err = s.signJWT(u)
	if err != nil {
		return "", "", err
	}

	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}
	refresh = base64.RawURLEncoding.EncodeToString(b)

	err = s.repo.CreateSession(ctx, Session{
		TokenHash: hashToken(refresh),
		UserID:    u.ID,
		ExpiresAt: time.Now().Add(s.refreshTTL),
	})
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (access, refresh string, err error) {
	u, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", "", ErrInvalidCredentials
	}

	return s.issueTokens(ctx, u)
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func (s *Service) Refresh(ctx context.Context, rawRefresh string) (access, newRefresh string, err error) {
	sess, err := s.repo.FindSession(ctx, hashToken(rawRefresh))
	if err != nil {
		return "", "", err
	}

	if time.Now().After(sess.ExpiresAt) {
		_ = s.repo.DeleteSession(ctx, sess.TokenHash)
		return "", "", ErrSessionInvalid
	}

	u, err := s.repo.FindUserByID(ctx, sess.UserID)
	if err != nil {
		return "", "", ErrSessionInvalid
	}

	if err := s.repo.DeleteSession(ctx, sess.TokenHash); err != nil {
		return "", "", err
	}

	return s.issueTokens(ctx, u)
}

func (s *Service) Logout(ctx context.Context, rawRefresh string) error {
	return s.repo.DeleteSession(ctx, hashToken(rawRefresh))
}
