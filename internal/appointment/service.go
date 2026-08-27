package appointment

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Appointment, error) {
	return s.repo.GetByID(ctx, id)
}
