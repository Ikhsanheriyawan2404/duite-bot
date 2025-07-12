package main

import (
	"log"
	"os"

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
	transactionHandler := handlers.NewTransactionHandler(conf.DB)
	categoryHandler := handlers.NewCategoryHandler(conf.DB)

	app.Get("/categories", middleware.JWTMiddleware, categoryHandler.GetCategories)

	app.Route("/transactions", func(r fiber.Router) {
		r.Use(middleware.JWTMiddleware)
	
		r.Post("/", transactionHandler.CreateTransaction)
		r.Get("/", transactionHandler.GetTransactions)
		r.Get("/:id", transactionHandler.GetTransactionByID)
		r.Put("/:id", transactionHandler.UpdateTransaction)
		r.Delete("/:id", transactionHandler.DeleteTransaction)
	})

	app.Post("/auth/magic-login", authHandler.MagicLogin)
	app.Post("/auth/refresh", authHandler.RefreshToken)
	app.Get("/user/me", middleware.JWTMiddleware, authHandler.GetMe)

	serverPort := os.Getenv("GATEWAY_PORT")
	
	log.Println("Gateway running on :" + serverPort)
	log.Fatal(app.Listen(":" + serverPort))
	
}
