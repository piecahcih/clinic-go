package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type postgresRepo struct {
	db *sqlx.DB
}

func NewPostgresRepo(db *sqlx.DB) *postgresRepo {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateUser(ctx context.Context, u *User) error {
	err := r.db.GetContext(ctx, u, `
		INSERT INTO users (id, first_name, last_name, birth_date, gender, email, password_hash, role)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id, first_name, last_name, birth_date, gender, email, password_hash, role, created_at`,
		u.ID, u.FirstName, u.LastName, u.BirthDate, u.Gender, u.Email, u.PasswordHash, u.Role)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrEmailTaken
		}
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *postgresRepo) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool

	err := r.db.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, email)
	if err != nil {
		return false, fmt.Errorf("check email exists %s: %w", email, err)
	}
	return exists, nil
}

func (r *postgresRepo) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	var u User

	err := r.db.GetContext(ctx, &u, `
	SELECT id, first_name, last_name, birth_date, gender, email, password_hash, role, created_at
	FROM users
	WHERE email = $1`, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("find user by email %s: %w", email, err)
	}
	return &u, nil
}

func (r *postgresRepo) FindUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var u User

	err := r.db.GetContext(ctx, &u, `
		SELECT id, first_name, last_name, birth_date, gender, email, password_hash, role, created_at
		FROM users
		WHERE id = $1`, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("find user by id %s: %w", id, err)
	}
	return &u, nil
}

func (r *postgresRepo) CreateSession(ctx context.Context, s Session) error {
	_, err := r.db.ExecContext(ctx, `
	INSERT INTO sessions (token, user_id, expires_at)
	VALUES ($1, $2, $3)`,
		s.TokenHash, s.UserID, s.ExpiresAt)
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindSession(ctx context.Context, hash string) (*Session, error) {
	var s Session
	err := r.db.GetContext(ctx, &s, `
		SELECT token, user_id, expires_at, created_at FROM sessions
		WHERE token = $1`, hash)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionInvalid
		}
		return nil, fmt.Errorf("find session: %w", err)
	}
	return &s, nil
}

func (r *postgresRepo) DeleteSession(ctx context.Context, hash string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM sessions WHERE token = $1`, hash)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

type Repository interface {
	CreateUser(ctx context.Context, u *User) error
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	FindUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	EmailExists(ctx context.Context, email string) (bool, error)

	CreateSession(ctx context.Context, s Session) error
	FindSession(ctx context.Context, tokenHash string) (*Session, error)
	DeleteSession(ctx context.Context, tokenHash string) error
}
