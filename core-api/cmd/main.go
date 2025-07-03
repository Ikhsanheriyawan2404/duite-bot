package main

import (
	"finance-bot/config"
	"finance-bot/handler"
	"finance-bot/repository"
	"finance-bot/service"
	"finance-bot/middleware"
	
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
)

const MaxLimitHit = 10;
var log = logrus.New()

func main() {
	appEnv := os.Getenv("APP_ENV")
	pathEnv := "../"
	if appEnv == "production" {
		pathEnv = "."
	}

	if err := config.LoadConfig(pathEnv); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	
	db := config.ConnectDB()
	config.InitRedis()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	transactionRepo := repository.NewTransactionRepository(db)
	transactionService := service.NewTransactionService(transactionRepo)

	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)

	// Create LLM service
	// aiService := service.NewLLMService(
	// 	service.ProviderDeepSeek,
	// 	config.AppConfig.LLMApiKey,
	// 	"https://api.deepseek.com/v1/chat/completions",
	// 	"deepseek-chat",
	// )
	aiService := service.NewLLMService(
		service.ProviderChatGPT,
		config.AppConfig.LLMApiKey,
		"",
		"gpt-4.1-nano",
	)
	
    userHandler := handler.NewUserHandler(userService, transactionService, categoryService, *aiService)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
        AllowOrigins: "*",
        AllowHeaders: "Origin, Content-Type, Accept",
    }))

	// Apply to all routes (1 requests per 1 seconds)
	app.Use(limiter.RateLimiterMiddleware(60, 60*time.Second))
		
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})
	
	users := app.Group("/users")
	
	users.Post("/register", userHandler.RegisterUser)
	users.Post("/magic-link", userHandler.GenerateMagicLink)

	users.Get("/:userId", userHandler.GetUser)
	users.Get("/:chatId/exists", userHandler.CheckUser)
	users.Get("/:userId/transactions", userHandler.GetTransactions)
	
	users.Get("/:chatId/transactions/daily", userHandler.GetDailyReport)
	users.Get("/:chatId/transactions/monthly", userHandler.GetMonthlyReport)
	users.Delete("/:chatId/transactions/:transactionId", userHandler.DeleteTransactionByID)
	
	users.Post("/:chatId/transactions/ai-classify", userHandler.ParseAndSaveTransaction)

	serverAddr := "0.0.0.0:" + config.AppConfig.ServerPort
	if err := app.Listen(serverAddr); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
