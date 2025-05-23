package main

import (
	"log"
	"os"

	"bot-tele/config"
	"bot-tele/handler"
	"bot-tele/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// load configuration
	appEnv := os.Getenv("APP_ENV")
	pathEnv := "../"
	if appEnv == "production" {
		pathEnv = "."
	}
	if err := config.LoadConfig(pathEnv); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// declare API client
	apiClient := service.NewAPIClient()

	// Initialize Telegram bot
	bot, err := tgbotapi.NewBotAPI(config.AppConfig.TelegramToken)
	if err != nil {
        log.Panic(err)
    }
    log.Println("Authorized on account " + bot.Self.UserName)
	
	u := tgbotapi.NewUpdate(0)
    u.Timeout = 60
    updates := bot.GetUpdatesChan(u)
	
    for update := range updates {

		if update.Message == nil {
			continue
		}

		handler.HandleCommandAndInput(update, bot, apiClient)
    }
}