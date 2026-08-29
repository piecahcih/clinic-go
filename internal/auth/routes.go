package auth

import (
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(r fiber.Router, h *Handler) {
	g := r.Group("/auth")

	g.Post("/register", h.Register)
	g.Post("/login", h.Login)
}
