package handlers

import "github.com/gofiber/fiber/v2"

func RetrieveUser(c *fiber.Ctx) error {
	username := c.Locals("username").(string)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"username": username})
}
