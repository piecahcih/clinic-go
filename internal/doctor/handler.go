package doctor

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) DoctorInfo(ctx fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid doctor id"})
	}

	// role, _ := ctx.Locals("role").(string)

	// doc, err := h.svc.DoctorInfo(ctx.Context(), id, role)
	doc, err := h.svc.DoctorInfo(ctx.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return ctx.Status(fiber.StatusNotFound).
				JSON(fiber.Map{"error": err.Error()})
		// case errors.Is(err, ErrNotADoctor):
		// 	return ctx.Status(fiber.StatusForbidden).
		// 		JSON(fiber.Map{"error": err.Error()})
		default:
			return ctx.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"error": "internal error"})
		}
	}

	return ctx.JSON(doc)
}

func (h *Handler) DoctorSlots(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid doctor id"})
	}

	dateStr := c.Query("date")
	day, err := time.Parse("2006-01-02", dateStr) //(layout,value)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "date must be YYYY-MM-DD"})
	}

	slots, err := h.svc.AvailableSlots(c.Context(), id, day)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal error"})
	}

	localSlots := make([]time.Time, len(slots))
	for i, t := range slots {
		localSlots[i] = t.In(bangkok)
	}

	return c.JSON(fiber.Map{"doctorId": id, "date": dateStr, "slots": localSlots})
}
