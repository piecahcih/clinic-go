package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func Authenticate(secret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		header := c.Get("Authorization")
		tokenStr, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "missing or malformed bearer token"})
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "invalid or expired token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "invalid token claims"})
		}

		sub, ok := claims["sub"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "invalid token claims"})
		}
		userID, err := uuid.Parse(sub)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "invalid token claims"})
		}
		role, _ := claims["role"].(string)

		c.Locals("patientID", userID)
		c.Locals("role", role)
		return c.Next()
	}
}

// // MockAuthCheck
// func MockAuthCheck() fiber.Handler {
// 	// mockID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
// 	// mockRole := "patient"
// 	mockID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
// 	mockRole := "doctor"

// 	return func(c fiber.Ctx) error {
// 		c.Locals("patientID", mockID)
// 		c.Locals("role", mockRole)
// 		return c.Next()
// 	}
// }
