package main

import (
	"go-miniblog/config"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()
	config.InitDB()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to Go Blog API"})
	})

	app.Get("/api/v1/auth/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Auth test endpoint works!",
		})
	})

	// Jalankan server di port dari .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port kalau belum ada di .env
	}

	log.Println("Server running on port:", port)
	app.Listen(":" + port)
}
