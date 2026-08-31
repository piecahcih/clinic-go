package doctor

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type postgresRepo struct {
	db *sqlx.DB
}

func NewPostgresRepo(db *sqlx.DB) *postgresRepo {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) DoctorInfo(ctx context.Context, id uuid.UUID) (*User, error) {
	var u User

	err := r.db.GetContext(ctx, &u, `
	SELECT id, first_name, last_name, role
	FROM users
	WHERE id = $1`, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get doctor info %s: %w", id, err)
	}
	return &u, nil
}

func (r *postgresRepo) BookedSlots(ctx context.Context, doctorID uuid.UUID, day time.Time) ([]time.Time, error) {
	var times []time.Time
	y, m, d := day.Date()
	start := time.Date(y, m, d, 0, 0, 0, 0, bangkok)
	end := start.Add(24 * time.Hour)

	err := r.db.SelectContext(ctx, &times, `
		SELECT start_time
		FROM appointments
		WHERE doctor_id = $1
		  AND status = 'booked'
		  AND start_time >= $2 AND start_time < $3`,
		doctorID, start, end)
	if err != nil {
		return nil, fmt.Errorf("booked slots %s: %w", doctorID, err)
	}
	return times, nil
}

func (r *postgresRepo) ListDoctors(ctx context.Context) ([]User, error) {
	var doctors []User
	err := r.db.SelectContext(ctx, &doctors, `
		SELECT id, first_name, last_name, role
		FROM users
		WHERE role = 'doctor'
		ORDER BY last_name, first_name`)
	if err != nil {
		return nil, fmt.Errorf("list doctors: %w", err)
	}
	return doctors, nil
}

type Repository interface {
	DoctorInfo(ctx context.Context, id uuid.UUID) (*User, error)
	BookedSlots(ctx context.Context, doctorID uuid.UUID, day time.Time) ([]time.Time, error)
	ListDoctors(ctx context.Context) ([]User, error)
}
