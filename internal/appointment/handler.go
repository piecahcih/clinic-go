package appointment

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type Handler struct {
	svc *Service
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

	appt, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return c.Status(fiber.StatusNotFound).
				JSON(fiber.Map{"error": "appointment not found"})
		}
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "internal error"})
	}

	return c.JSON(appt)
}
