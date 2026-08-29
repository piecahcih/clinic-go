package doctor

import (
	"github.com/gofiber/fiber/v3"
	"github.com/piecahcih/clinic-go/internal/middleware"
)

func RegisterRoutes(r fiber.Router, h *Handler) {
	g := r.Group("/doctor")

	// doctor details
	g.Get("/:id", middleware.MockAuthCheck(), h.DoctorInfo)
	// available booking times
	g.Get("/:id/slots", h.DoctorSlots)
}
