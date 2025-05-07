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

type ResetPasswordRequest struct {
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	Username     string `json:"username"`
	RefreshToken string `json:"refreshToken"`
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

	refreshToken, err := auth.GenerateRefreshToken(req.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate refresh token."})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token, "refreshToken": refreshToken})
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

	refreshToken, err := auth.GenerateRefreshToken(user.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate refresh token."})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token, "refreshToken": refreshToken})
}

func Logout(c *fiber.Ctx) error {
	var req RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request."})
	}
	err := database.InMemoryDB().DeleteRefreshToken(req.Username, req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete refresh token."})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logout successful."})
}

func ResetPassword(c *fiber.Ctx) error {
	var req ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request."})
	}

	user, err := database.GetUserByUsername(c.Locals("username").(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not get user."})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not hash password."})
	}

	// Update user password
	user.Password = string(hashedPassword)

	err = database.UpdateUser(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update user's password."})
	}

	// Delete refresh tokens
	err = database.InMemoryDB().DeleteRefreshTokens(user.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete refresh tokens."})
	}

	// Generate new refresh token
	refreshToken, err := auth.GenerateRefreshToken(user.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate refresh token."})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Password reset successfully.", "refreshToken": refreshToken})
}

func RefreshToken(c *fiber.Ctx) error {
	var req RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request."})
	}

	userRefreshTokens, err := database.InMemoryDB().GetRefreshTokens(req.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not get refresh tokens."})
	}

	for _, rToken := range userRefreshTokens {
		if rToken == req.RefreshToken {
			token, err := auth.GenerateJWT(req.Username)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token."})
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
		}
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token."})
}
