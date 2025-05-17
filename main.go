package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"finance-bot/config"
	"finance-bot/handler"
	"finance-bot/model"
	"finance-bot/repository"
	"finance-bot/service"
	"finance-bot/utils"

	telegrambot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const MaxLimitHit = 10;

var db *gorm.DB
var log = logrus.New()

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model string `json:"model"`
	Messages []ChatMessage `json:"messages"`
	MaxTokens int `json:"max_tokens,omitempty"`
	ResponseFormat map[string]any `json:"response_format,omitempty"`
}

type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`

	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type Result struct {
	TransactionType string `json:"type"`
	Amount float64 `json:"amount"`
	Category string `json:"category"`
}

func classifyTransaction(desc string) (*Result, *ChatResponse, error) {
	prompt := fmt.Sprintf(utils.PromptDefault, desc)

	requestBody := ChatRequest{
		Model: "deepseek-chat",
		Messages: []ChatMessage{
			{Role: "user", Content: prompt},
		},
		ResponseFormat: map[string]any{ 
			"type": "json_object",
		},
		// MaxTokens: 50,
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", config.AppConfig.LLMApiUrl, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + config.AppConfig.LLMApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var chatResp ChatResponse
	json.Unmarshal(body, &chatResp)
	jsonStr := chatResp.Choices[0].Message.Content

	var result Result
	json.Unmarshal([]byte(jsonStr), &result)

	return &result, &chatResp, err
}

func connectDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.AppConfig.DBHost,
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBName,
		config.AppConfig.DBPort,
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("Failed to connect DB:", err)
	}

	db.AutoMigrate(&model.Transaction{})
	db.AutoMigrate(&model.User{})
}

func initializeRestApi() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
        AllowOrigins: "*",
        AllowHeaders: "Origin, Content-Type, Accept",
    }))

	app.Static("/", "./public")
	
	api := app.Group("/api")
	users := api.Group("/users")

	userRepo := repository.NewUserRepository(db)
    userService := service.NewUserService(userRepo)
    userHandler := handler.NewUserHandler(userService)

	users.Get("/:userId", userHandler.GetUser)
	users.Get("/:userId/transactions", userHandler.GetTransactions)
	
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	go func() {
		serverAddr := "0.0.0.0:" + config.AppConfig.ServerPort
		if err := app.Listen(serverAddr); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()
}

func main() {
	// Coba buka file log
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// Gunakan stdout sebagai fallback
		log.SetOutput(os.Stdout)
		log.WithError(err).Warn("Gagal membuka file log, menggunakan stdout")
	} else {
		log.SetOutput(file)
	}
	
	if err := config.LoadConfig("."); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	connectDB()
	initializeRestApi()

	bot, err := telegrambot.NewBotAPI(config.AppConfig.TelegramToken)
	if err != nil {
        log.Panic(err)
    }

	bot.Debug = true
    log.Info("Authorized on account " + bot.Self.UserName)

	u := telegrambot.NewUpdate(0)
    u.Timeout = 60

    updates := bot.GetUpdatesChan(u)

    for update := range updates {
		start := time.Now() // mulai timer
		
        if update.Message == nil {
			continue
		}
		
		chatID := update.Message.Chat.ID
		text := update.Message.Text
		parts := strings.Fields(text)
		command := parts[0]

		// Kondisi untuk command
		switch command {
		case "/start":
			helpText := utils.HelpText
			msg := telegrambot.NewMessage(chatID, helpText)
			msg.ParseMode = "Markdown"
			bot.Send(msg)
			continue;
		case "/harian":
			transactions, _ := service.GetDailyReport(db, chatID)
			msgText := service.FormatDailyReport(transactions)
			report := msgText
			msg := telegrambot.NewMessage(chatID, utils.EscapeMarkdown(report))
			msg.ParseMode = "MarkdownV2"
			bot.Send(msg)
			continue;
		case "/bulanan":
			transactions, _ := service.GetDailyReport(db, chatID)
			msgText := service.FormatMonthlyReport(transactions)
			report := msgText
			msg := telegrambot.NewMessage(chatID, utils.EscapeMarkdown(report))
			msg.ParseMode = "MarkdownV2"
			bot.Send(msg)
			continue;
		case "/hapus":
			parts := strings.Fields(text)
			if len(parts) < 2 {
				msg := telegrambot.NewMessage(chatID, "Format salah. Gunakan: /hapus {ID transaksi}")
				bot.Send(msg)
				continue
			}
			
			indexStr := parts[1]
			transactionId, err := strconv.Atoi(indexStr)
			if err != nil {
				msg := telegrambot.NewMessage(chatID, "ID transaksi harus berupa angka.")
				bot.Send(msg)
				continue
			}

			err = service.DeleteTransactionByID(db, uint(transactionId), chatID)
			if (err != nil) {
				msg := telegrambot.NewMessage(chatID, "Transaksi tidak ditemukan")
				bot.Send(msg)
				continue
			}
			msg := telegrambot.NewMessage(chatID, "Transaksi berhasil dihapus!")
			bot.Send(msg)
			continue;
		case "/daftar":
			parts := strings.Fields(text)
			if len(parts) < 2 {
				msg := telegrambot.NewMessage(chatID, "Format salah. Gunakan: /daftar {nama}")
				bot.Send(msg)
				continue
			}
			userFullName := parts[1]

			user, _ := service.RegisterUser(db, chatID, userFullName)
			msg := telegrambot.NewMessage(chatID, "Selamat datang " + user.Name)
			bot.Send(msg)
			continue;
		case "/dashboard":
			userExist := service.CheckUser(db, chatID)
			if (!userExist) {
				msg := telegrambot.NewMessage(chatID, "Mohon melakukan daftar terlebih dahulu")
				bot.Send(msg)
				continue
			}
			msg := telegrambot.NewMessage(chatID, "Klik tombol di bawah untuk membuka dashboard:")
			msg.ReplyMarkup = telegrambot.InlineKeyboardMarkup{
				InlineKeyboard: [][]telegrambot.InlineKeyboardButton{
					{
						telegrambot.NewInlineKeyboardButtonURL("📊 Buka Dashboard",
							config.AppConfig.DashboardUrl+"?ref="+utils.EncodeChatID(chatID)),
					},
				},
			}
			bot.Send(msg)
			continue;
		default:
			if strings.HasPrefix(text, "/") {
				msg := telegrambot.NewMessage(chatID, "❓ Command tidak dikenali. Gunakan /help untuk melihat daftar command.")
				bot.Send(msg)
				continue;
			}
		}

		// Cek otoritas user
		countTransactionUser, _ := service.CountTransactionsById(db, chatID)
		if (MaxLimitHit <= countTransactionUser) {
			msg := telegrambot.NewMessage(chatID, "Maaf anda sudah melebihi limit, hubungi admin di https://t.me/Hirumakun.")
			bot.Send(msg)
			continue;
		}

		// Klasifikasi prompt & buat transaksi
		result, fullResponse, _ := classifyTransaction(update.Message.Text)

		transactionType, err := utils.ParseTransactionType(result.TransactionType)
		if err != nil {
			// msgTxt := ""
			msg := telegrambot.NewMessage(chatID, "Transaksi gagal dibuat!")
			bot.Send(msg)
		}

		tx := model.Transaction{
			ChatID: chatID,
			OriginalText: update.Message.Text,
			TransactionType: transactionType,
			Amount: result.Amount,
			Category: utils.Slugify(result.Category),
			TransactionDate: time.Now(),
		}
		
		db.Create(&tx)

		// Balas Konfirmasi
		transactionTypeIndo, _ := utils.FormatTransactionType(result.TransactionType)

		msgText := fmt.Sprintf(
			"✅ Transaksi berhasil dicatat!\n\n📂 Tipe: %s\n💰 Jumlah: %s\n🏷️ Kategori: %s\n",
			transactionTypeIndo,
			utils.FormatRupiah(result.Amount),
			result.Category,
		)

		msg := telegrambot.NewMessage(chatID, msgText)
		bot.Send(msg)

		elapsed := time.Since(start)

		// Logging akhir
		log.Info("Total token " + strconv.Itoa(fullResponse.Usage.TotalTokens))
		log.Info("Execution code time ", elapsed)
    }
}
