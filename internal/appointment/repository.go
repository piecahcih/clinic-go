package appointment

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type postgresRepo struct {
	db *sqlx.DB
}

func NewPostgresRepo(db *sqlx.DB) *postgresRepo {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) GetByID(ctx context.Context, id uuid.UUID) (*Appointment, error) {
	var a Appointment

	err := r.db.GetContext(ctx, &a, `
		SELECT id, patient_id, doctor_id, start_time, description, status, created_at
		FROM appointments
		WHERE id = $1`, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get appointment %s: %w", id, err)
	}
	return &a, nil
}

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Appointment, error)
}
