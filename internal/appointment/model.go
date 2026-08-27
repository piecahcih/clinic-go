package appointment

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Appointment struct {
	ID          uuid.UUID
	PatientID   uuid.UUID
	DoctorID    uuid.UUID
	StartTime   time.Time
	Description *string
	Status      string
	CreatedAt   time.Time
}

var ErrNotFound = errors.New("appointment not found")
