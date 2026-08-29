package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

type Repository interface {
	CreateUser(ctx context.Context, u *User) error
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	EmailExists(ctx context.Context, email string) (bool, error)
}
