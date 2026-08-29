package appointment

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

func (r *postgresRepo) AllMyAppointment(ctx context.Context, userID uuid.UUID, role string) ([]Appointment, error) {
	var appts []Appointment

	column := "patient_id"
	if role == "doctor" {
		column = "doctor_id"
	}

	query := fmt.Sprintf(`
		SELECT id, patient_id, doctor_id, start_time, description, status, created_at
		FROM appointments
		WHERE %s = $1
		ORDER BY start_time`, column)

	err := r.db.SelectContext(ctx, &appts, query, userID)

	if err != nil {
		return nil, fmt.Errorf("list appointments: %w", err)
	}
	return appts, nil
}

func (r *postgresRepo) AddAppointment(ctx context.Context, a Appointment) (*Appointment, error) {
	a.ID = uuid.New()

	err := r.db.GetContext(ctx, &a, `
		INSERT INTO appointments (id, patient_id, doctor_id, start_time, description, status)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id, patient_id, doctor_id, start_time, description, status, created_at`,
		a.ID, a.PatientID, a.DoctorID, a.StartTime, a.Description, a.Status)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23P01" {
			return nil, ErrSlotConflict
		}
		return nil, fmt.Errorf("add appointment: %w", err)
	}
	return &a, nil
}

func (r *postgresRepo) CancelAppointment(ctx context.Context, id uuid.UUID) (*Appointment, error) {
	var a Appointment
	err := r.db.GetContext(ctx, &a, `
		UPDATE appointments
		SET status = 'cancelled'
		WHERE id = $1
		RETURNING id, patient_id, doctor_id, start_time, description, status, created_at`,
		id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("cancel appointment %s: %w", id, err)
	}
	return &a, nil
}

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Appointment, error)
	AllMyAppointment(ctx context.Context, userID uuid.UUID, role string) ([]Appointment, error)
	AddAppointment(ctx context.Context, a Appointment) (*Appointment, error)
	CancelAppointment(ctx context.Context, id uuid.UUID) (*Appointment, error)
}
