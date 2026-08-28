package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/piecahcih/clinic-go/internal/httperr"
)

func NotFound() fiber.Handler {
	return func(c fiber.Ctx) error {
		return httperr.NotFound("route " + c.Path() + " doesn't exist")
	}
}
