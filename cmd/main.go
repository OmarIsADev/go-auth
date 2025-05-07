package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/omarisadev/go-auth/database"
	"github.com/omarisadev/go-auth/handlers"
	"github.com/omarisadev/go-auth/middleware"
)

func main() {
	database.DBConnect()

	app := fiber.New()
	// public routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Post("/regester", handlers.Regester)
	app.Post("/login", handlers.Login)
	app.Post("/logout", handlers.Logout)
	app.Post("/refresh-token", handlers.RefreshToken)

	// protected routes
	app.Use(middleware.Protected())
	app.Get("/user", handlers.RetrieveUser)
	app.Post("/reset-password", handlers.ResetPassword)

	app.Listen(":3000")
}
