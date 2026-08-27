package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/piecahcih/clinic-go/internal/appointment"
)

func main() {
	app := fiber.New()
	api := app.Group("/api/v1")

	apptRepo := appointment.NewMemoryRepo()
	apptSvc := appointment.NewService(apptRepo)
	apptH := appointment.NewHandler(apptSvc)

	appointment.RegisterRoutes(api, apptH)

	// app.Get("/", func(c fiber.Ctx) error {
	// 	return c.SendString("test")
	// })

	app.Listen(":3000")
}
