package doctor

import (
	"errors"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	FirstName string    `db:"first_name" json:"firstName"`
	LastName  string    `db:"last_name" json:"lastName"`
	Role      string    `db:"role" json:"role"`
}

var ErrNotFound = errors.New("doctor not found")
var ErrNotADoctor = errors.New("You're not a doctor")
