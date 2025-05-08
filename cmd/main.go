package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Graceful shutdown handling
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down server...")
		log.Println("Saving memory DB...")

		err := database.InMemoryDB().Save()
		if err != nil {
			log.Fatalf("Failed to save memory DB: %v", err)
		}
		log.Println("Memory DB saved.")

		// Optional: Add a timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := app.ShutdownWithContext(ctx); err != nil {
			log.Fatalf("Server shutdown error: %v", err)
		}
		log.Println("Server gracefully shut down.")
	}()

	app.Listen(":3000")
}
