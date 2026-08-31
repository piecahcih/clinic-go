package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/joho/godotenv"
	"github.com/piecahcih/clinic-go/internal/appointment"
	"github.com/piecahcih/clinic-go/internal/auth"
	"github.com/piecahcih/clinic-go/internal/db"
	"github.com/piecahcih/clinic-go/internal/doctor"
	"github.com/piecahcih/clinic-go/internal/httperr"
	"github.com/piecahcih/clinic-go/internal/middleware"
)

func main() {
	_ = godotenv.Load()

	database, err := db.Connect(os.Getenv("DATABASE_URL"))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
	app.Use(recover.New()) //panic recover
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowCredentials: true,
	}))

	api := app.Group("/api/v1")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}
	jwtTTLMinutes, err := strconv.Atoi(os.Getenv("JWT_TTL_MINUTES"))
	if err != nil {
		log.Fatalf("JWT_TTL_MINUTES: %v", err)
	}
	refreshTTLDays, err := strconv.Atoi(os.Getenv("REFRESH_TTL_DAYS"))
	if err != nil {
		log.Fatalf("REFRESH_TTL_DAYS: %v", err)
	}
	refreshTTL := time.Duration(refreshTTLDays) * 24 * time.Hour
	secureCookies, _ := strconv.ParseBool(os.Getenv("COOKIE_SECURE")) // defaults false, e.g. for local http

	authRepo := auth.NewPostgresRepo(database)
	authSvc := auth.NewService(authRepo, jwtSecret, time.Duration(jwtTTLMinutes)*time.Minute, refreshTTL)
	authH := auth.NewHandle(authSvc, refreshTTL, secureCookies)

	auth.RegisterRoutes(api, authH)

	authMW := middleware.Authenticate(jwtSecret)

	// apptRepo := appointment.NewMemoryRepo()
	apptRepo := appointment.NewPostgresRepo(database)
	apptSvc := appointment.NewService(apptRepo)
	go appointment.RunNoShowWorker(ctx, apptRepo, 15*time.Minute)
	apptH := appointment.NewHandler(apptSvc)

	appointment.RegisterRoutes(api, apptH, authMW)

	docRepo := doctor.NewPostgresRepo(database)
	docSvc := doctor.NewService(docRepo)
	docH := doctor.NewHandler(docSvc)

	doctor.RegisterRoutes(api, docH, authMW)

	app.Use(middleware.NotFound())

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// The thing we just add is a handler foe server quit, error.
	// We added a handler so that when the OS sends a termination signal
	// whether from a developer's Ctrl+C or a deploy platform's shutdown request
	// our program gets a chance to react instead of dying instantly.

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- app.Listen(fmt.Sprintf(":%s", port))
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		log.Fatalf("server error: %v", err)
	case sig := <-quit:
		log.Printf("received %v, starting graceful shutdown", sig)
	}

	// select waits on two channels: if the server itself fails to start, we log the error and exit hard (nothing to gracefully close, since it never opened).
	// If a termination signal arrives instead, we log that we're shutting down gracefully and give in-flight requests up to 10 seconds to finish before the process exits."

	cancel()

	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		log.Printf("graceful shutdown error: %v", err)
	}
}
