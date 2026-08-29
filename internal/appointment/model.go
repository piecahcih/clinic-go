package appointment

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Appointment struct {
	ID          uuid.UUID `db:"id" json:"id"`
	PatientID   uuid.UUID `db:"patient_id" json:"patientId"`
	DoctorID    uuid.UUID `db:"doctor_id" json:"doctorId"`
	StartTime   time.Time `db:"start_time" json:"startTime"`
	Description *string   `db:"description" json:"description"`
	Status      string    `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
}

var ErrNotFound = errors.New("appointment not found")
var ErrSlotConflict = errors.New("slot conflicts, you already book on this time or doctor are booked")
var ErrForbidden = errors.New("you can only cancel your own appointment")
var ErrInvalidState = errors.New("appointment is not in a cancellable state")
var ErrCancelWindowPassed = errors.New("appointments can only be cancelled at least 24 hours before the start time")

// Layer 	 |	Convention	| Example
//-------------------------------------------
// SQL column| snake_case	| patient_id
// Go struct | PascalCase	| PatientID
// JSON      | camelCase	| patientId
