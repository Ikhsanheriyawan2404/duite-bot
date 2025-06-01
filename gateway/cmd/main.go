package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	
	"gateway/config"
	"gateway/handlers"
	"gateway/middleware"
)

func main() {
	conf := config.Load()

	app := fiber.New()
	app.Use(cors.New())

	authHandler := handlers.NewAuthHandler(conf.DB, conf.Redis)

	app.Post("/auth/magic-login", authHandler.MagicLogin)
	app.Post("/auth/refresh", authHandler.RefreshToken)
	app.Get("/user/me", middleware.JWTMiddleware, authHandler.GetMe)

	log.Println("Gateway running on :8080")
	log.Fatal(app.Listen(":8080"))
}
