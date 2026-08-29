package auth

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	svc *Service
}

type RegisterDTO struct {
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	BirthDate time.Time `json:"birthDate"`
	Gender    string    `json:"gender"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
}

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewHandle(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(c fiber.Ctx) error {
	var req RegisterDTO

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.Password == "" ||
		req.Gender == "" || req.BirthDate.IsZero() {
		return c.Status(fiber.StatusUnprocessableEntity).
			JSON(fiber.Map{"error": "firstName, lastName, birthDate, gender, email and password are required"})
	}

	u, err := h.svc.Register(c.Context(), User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		BirthDate: req.BirthDate,
		Gender:    req.Gender,
		Email:     req.Email,
		Role:      req.Role,
	}, req.Password)

	if err != nil {
		switch {
		case errors.Is(err, ErrEmailTaken):
			return c.Status(fiber.StatusConflict).
				JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"error": "internal error"})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(u)
}

func (h *Handler) Login(c fiber.Ctx) error {
	var req LoginDTO

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusUnprocessableEntity).
			JSON(fiber.Map{"error": "email and password are required"})
	}

	token, err := h.svc.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"error": "internal error"})
		}
	}

	return c.JSON(fiber.Map{"token": token})
}
