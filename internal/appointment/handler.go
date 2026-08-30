package appointment

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type Handler struct {
	svc *Service
}

type createAppointmentDTO struct {
	DoctorID    uuid.UUID `json:"doctorId"`
	StartTime   time.Time `json:"startTime"`
	Description *string   `json:"description"`
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetByID(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid appointment id"})
	}

	callerID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"error": "unauthorized"})
	}

	appt, err := h.svc.GetByID(c.Context(), id, callerID)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return c.Status(fiber.StatusNotFound).
				JSON(fiber.Map{"error": "appointment not found"})
		case errors.Is(err, ErrForbidden):
			return c.Status(fiber.StatusForbidden).
				JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"error": "internal error"})
		}
	}

	return c.JSON(appt)
}

func (h *Handler) AllMyAppointment(c fiber.Ctx) error {
	callerID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"error": "unauthorized"})
	}
	role, _ := c.Locals("role").(string)

	appts, err := h.svc.AllMyAppointment(c.Context(), callerID, role)
	if err != nil {
		return err
	}
	return c.JSON(appts)
}

func (h *Handler) AddAppointment(c fiber.Ctx) error {
	var req createAppointmentDTO

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Invalid request body"})
	}

	callerID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"error": "unauthorized"})
	}
	role, _ := c.Locals("role").(string)

	if req.DoctorID == uuid.Nil || req.StartTime.IsZero() {
		return c.Status(fiber.StatusUnprocessableEntity).
			JSON(fiber.Map{"error": "doctor_id and start_time are required"})
	}

	appt, err := h.svc.AddAppointment(c.Context(), Appointment{
		PatientID:   callerID,
		DoctorID:    req.DoctorID,
		StartTime:   req.StartTime,
		Description: req.Description,
		Status:      "booked",
	}, role)

	if err != nil {
		switch {
		case errors.Is(err, ErrSlotConflict):
			return c.Status(fiber.StatusConflict).
				JSON(fiber.Map{"error": err.Error()})
		case errors.Is(err, ErrForbidden):
			return c.Status(fiber.StatusForbidden).
				JSON(fiber.Map{"error": "only patients can book appointments"})
		default:
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"error": "internal error"})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(appt)
}

func (h *Handler) CancelAppointment(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid appointment id"})
	}

	callerID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"error": "unauthorized"})
	}

	appt, err := h.svc.CancelAppointment(c.Context(), id, callerID)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		case errors.Is(err, ErrForbidden):
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		case errors.Is(err, ErrInvalidState), errors.Is(err, ErrCancelWindowPassed):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal error"})
		}
	}
	return c.JSON(appt)
}
