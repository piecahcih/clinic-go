package appointment

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Appointment, error)
}

// mock up data
type memoryRepo struct {
	items map[uuid.UUID]Appointment
}

func NewMemoryRepo() *memoryRepo {
	seed := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	return &memoryRepo{
		items: map[uuid.UUID]Appointment{
			seed: {
				ID:        seed,
				PatientID: uuid.MustParse("22222222-2222-2222-2222-222222222222"),
				DoctorID:  uuid.MustParse("33333333-3333-3333-3333-333333333333"),
				StartTime: time.Date(2026, 9, 1, 9, 0, 0, 0, time.UTC),
				Status:    "booked",
				CreatedAt: time.Now().UTC(),
			},
		},
	}
}

func (r *memoryRepo) GetByID(_ context.Context, id uuid.UUID) (*Appointment, error) {
	a, ok := r.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &a, nil
}
