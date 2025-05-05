package middleware

import (
	"github.com/omarisadev/go-auth/auth" // Adjust import path if needed
	"github.com/omarisadev/go-auth/database"

	"github.com/gofiber/fiber/v2"
)

// Protected is a Fiber middleware to check for a valid JWT.
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Missing authorization header"})
		}

		token := ""
		// Assuming Bearer token format: "Bearer <token>"
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid authorization format"})
		}

		username, err := auth.VerifyJWT(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token"})
		}

		_, err = database.GetUserByUsername(username)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
		}

		// Optionally, you can store the username in the Fiber context for later use
		c.Locals("username", username)
		return c.Next()
	}
}
