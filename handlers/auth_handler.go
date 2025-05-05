package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/omarisadev/go-auth/auth"
	"github.com/omarisadev/go-auth/database"
	"github.com/omarisadev/go-auth/models"
	"golang.org/x/crypto/bcrypt"
)

type RegesterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Regester(c *fiber.Ctx) error {
	var req RegesterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request."})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not hash password."})
	}

	if err := database.CreateUser(&models.User{
		Username: req.Username,
		Password: string(hashedPassword),
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create user."})
	}

	token, err := auth.GenerateJWT(req.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token."})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request."})
	}

	user, err := database.GetUserByUsername(req.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not get user."})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials."})
	}

	token, err := auth.GenerateJWT(user.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token."})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}
