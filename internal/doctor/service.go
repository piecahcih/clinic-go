package doctor

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

func (s *Service) DoctorInfo(ctx context.Context, id uuid.UUID, role string) (*User, error) {
	if role != "doctor" {
		return nil, ErrNotADoctor
	}
	u, err := s.repo.DoctorInfo(ctx, id)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Thailand has no DST, so a fixed offset is safe and avoids depending on the
// system having the IANA tzdata (Asia/Bangkok) installed.
var bangkok = time.FixedZone("+07:00", 7*60*60)

const (
	workDayStart = 8 * time.Hour
	workDayEnd   = 18 * time.Hour
	slotDuration = 15 * time.Minute
)

func (s *Service) AvailableSlots(ctx context.Context, doctorID uuid.UUID, day time.Time) ([]time.Time, error) {
	booked, err := s.repo.BookedSlots(ctx, doctorID, day)
	if err != nil {
		return nil, err
	}
	bookedSet := make(map[time.Time]struct{}, len(booked))
	for _, t := range booked {
		bookedSet[t.UTC()] = struct{}{}
	}

	y, m, d := day.Date()
	dayStart := time.Date(y, m, d, 0, 0, 0, 0, bangkok)

	slots := make([]time.Time, 0)
	for t := dayStart.Add(workDayStart); t.Before(dayStart.Add(workDayEnd)); t = t.Add(slotDuration) {
		if _, taken := bookedSet[t.UTC()]; taken {
			continue
		}
		if t.Before(time.Now()) { // drop past slots
			continue
		}
		slots = append(slots, t)
	}
	return slots, nil
}
