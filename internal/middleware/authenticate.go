package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// MockAuthCheck
func MockAuthCheck() fiber.Handler {
	// mockID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	// mockRole := "patient"
	mockID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	mockRole := "doctor"

	return func(c fiber.Ctx) error {
		c.Locals("patientID", mockID)
		c.Locals("role", mockRole)
		return c.Next()
	}
}
