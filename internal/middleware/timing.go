package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
)

// measures how long a request takes.
func Timing() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		log.Printf("%s %s %v", c.Method(), c.Path(), time.Since(start))
		return err
	}
}
