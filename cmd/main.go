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

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Post("/regester", handlers.Regester)
	app.Post("/login", handlers.Login)

	// protected routes
	app.Get("/user", middleware.Protected(), handlers.RetrieveUser)

	app.Listen(":3000")
}
