package appointment

import (
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(r fiber.Router, h *Handler, authMW fiber.Handler) {
	g := r.Group("/appointments")

	g.Get("/", authMW, h.AllMyAppointment)
	g.Get("/:id", authMW, h.GetByID)
	g.Post("/", authMW, h.AddAppointment)
	g.Patch("/:id", authMW, h.CancelAppointment)
}
