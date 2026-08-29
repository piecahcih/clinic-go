package appointment

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Appointment, error) {
	appt, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if appt.PatientID != userID {
		return nil, ErrForbidden
	}
	return appt, nil
}

func (s *Service) AllMyAppointment(ctx context.Context, userID uuid.UUID, role string) ([]Appointment, error) {
	return s.repo.AllMyAppointment(ctx, userID, role)
}

func (s *Service) AddAppointment(ctx context.Context, a Appointment, role string) (*Appointment, error) {
	if role != "patient" {
		return nil, ErrForbidden
	}
	return s.repo.AddAppointment(ctx, a)
}

func (s *Service) CancelAppointment(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Appointment, error) {
	appt, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if appt.PatientID != userID {
		return nil, ErrForbidden
	}
	if appt.Status != "booked" {
		return nil, ErrInvalidState
	}
	if time.Until(appt.StartTime) < 24*time.Hour {
		return nil, ErrCancelWindowPassed
	}

	return s.repo.CancelAppointment(ctx, id)
}
