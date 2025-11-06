package main

import (
	"fmt"
	"go-miniblog/config"
	"go-miniblog/internal/auth"
	"go-miniblog/internal/middleware"
	"go-miniblog/internal/post"

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
	if err := config.DB.AutoMigrate(&auth.User{}, &post.Post{}); err != nil {
		log.Fatal("❌ Migration failed:", err)
	}
	fmt.Println("✅ Database migrated")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to Go Blog API"})
	})

	app.Get("/api/v1/auth/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Auth test endpoint works!",
		})
	})

	// Auth Routes (Public)
	api := app.Group("/api")
	api.Post("/register", auth.RegisterHandler)
	api.Post("/login", auth.LoginHandler)

	// Post Routes (Public - Read Only)
	api.Get("/posts", post.GetAllPostsHandler)
	api.Get("/posts/:id", post.GetPostByIDHandler)
	api.Get("/users/:user_id/posts", post.GetPostsByUserHandler)

	// Protected Routes (Require Authentication)
	protected := api.Group("", middleware.AuthRequired())
	protected.Post("/posts", post.CreatePostHandler)
	protected.Put("/posts/:id", post.UpdatePostHandler)
	protected.Delete("/posts/:id", post.DeletePostHandler)

	// Jalankan server di port dari .env
	port := os.Getenv("PORT")

	log.Println("Server running on port:", port)
	app.Listen(":" + port)
}
