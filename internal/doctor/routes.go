package doctor

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(r fiber.Router) {
	g := r.Group("/doctor")

	// doctor details
	g.Get("/:id", func(c fiber.Ctx) error {
		return c.SendString("GET /doctor/:id")
	})
	// available booking times
	g.Get("/:id/slots", func(c fiber.Ctx) error {
		return c.SendString("GET /doctor/:id/slots")
	})
}
