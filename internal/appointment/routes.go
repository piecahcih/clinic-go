package appointment

import (
	"github.com/gofiber/fiber/v3"
	"github.com/piecahcih/clinic-go/internal/middleware"
)

func RegisterRoutes(r fiber.Router, h *Handler) {
	g := r.Group("/appointments")

	g.Get("/", middleware.MockAuthCheck(), h.AllMyAppointment)
	g.Get("/:id", middleware.MockAuthCheck(), h.GetByID)
	g.Post("/", middleware.MockAuthCheck(), h.AddAppointment)
	g.Patch("/:id", middleware.MockAuthCheck(), h.CancelAppointment)
}
