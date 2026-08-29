package main

import (
	"errors"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/joho/godotenv"
	"github.com/piecahcih/clinic-go/internal/appointment"
	"github.com/piecahcih/clinic-go/internal/db"
	"github.com/piecahcih/clinic-go/internal/doctor"
	"github.com/piecahcih/clinic-go/internal/httperr"
	"github.com/piecahcih/clinic-go/internal/middleware"
)

func main() {
	_ = godotenv.Load()

	database, err := db.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close() //run when main return ja
	log.Println("database connnected")

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			var ae *httperr.APIError
			if errors.As(err, &ae) {
				return c.Status(ae.Status).JSON(fiber.Map{"error": ae})
			}

			log.Printf("unhandled error: %v", err)
			return c.Status(500).JSON(fiber.Map{
				"error": fiber.Map{"code": "internal", "message": "internal server error"},
			})
		},
	})

	app.Use(requestid.New())
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	api := app.Group("/api/v1")

	// apptRepo := appointment.NewMemoryRepo()
	apptRepo := appointment.NewPostgresRepo(database)
	apptSvc := appointment.NewService(apptRepo)
	apptH := appointment.NewHandler(apptSvc)

	appointment.RegisterRoutes(api, apptH)

	docRepo := doctor.NewPostgresRepo(database)
	docSvc := doctor.NewService(docRepo)
	docH := doctor.NewHandler(docSvc)

	doctor.RegisterRoutes(api, docH)

	// app.Get("/", func(c fiber.Ctx) error {
	// 	return c.SendString("test")
	// })

	app.Use(middleware.NotFound())

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":3000"))
}
