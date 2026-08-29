package doctor

import (
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(r fiber.Router, h *Handler, authMW fiber.Handler) {
	g := r.Group("/doctor")

	// doctor details
	g.Get("/:id", authMW, h.DoctorInfo)
	// available booking times
	g.Get("/:id/slots", h.DoctorSlots)
}
