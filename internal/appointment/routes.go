package appointment

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(r fiber.Router, h *Handler) {
	g := r.Group("/appointments")

	g.Get("/:id", h.GetByID)
}
