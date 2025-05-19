package main

import (
	"bytes"
	"context"
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
	"finance-bot/internal/appcontext"
	"finance-bot/model"
	"finance-bot/repository"
	"finance-bot/service"
	"finance-bot/utils"

	telegrambot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

const MaxLimitHit = 10;

var log = logrus.New()
var appCtx *appcontext.AppContext

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

func hitChatGpt(desc string) (*Result, *ChatResponse, error) {

	prompt := fmt.Sprintf(utils.PromptDefault, desc)

	client := openai.NewClient(
		option.WithAPIKey(config.AppConfig.LLMApiKey), // Gantilah dengan API key Anda
	)

	chatCompletion, err := client.Chat.Completions.New(
		context.TODO(),
		openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(prompt),
			},
			Model: "gpt-4.1-nano",
			ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
				OfJSONObject: &shared.ResponseFormatJSONObjectParam{
					Type: "json_object",
				},
			},
		},
	)
	if err != nil {
		panic(err.Error())
	}
	jsonBytes, err := json.MarshalIndent(chatCompletion, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	var chatResp ChatResponse
	json.Unmarshal(jsonBytes, &chatResp)
	jsonStr := chatResp.Choices[0].Message.Content

	var result Result
	json.Unmarshal([]byte(jsonStr), &result)

	return &result, &chatResp, nil
}



func initializeRestApi(appCtx *appcontext.AppContext) {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
        AllowOrigins: "*",
        AllowHeaders: "Origin, Content-Type, Accept",
    }))

	app.Static("/", "./public")
	
	api := app.Group("/api")
	users := api.Group("/users")

	
    userHandler := handler.NewUserHandler(appCtx.UserService)

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

	db := config.ConnectDB()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	transactionRepo := repository.NewTransactionRepository(db)
	transactionService := service.NewTransactionService(transactionRepo)

	appCtx = &appcontext.AppContext{
		UserService:        userService,
		TransactionService: transactionService,
		UserStateStore:    make(map[int64]string),
	}

	initializeRestApi(appCtx)
	
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
		textMessage := update.Message.Text
		parts := strings.Fields(textMessage)
		command := parts[0]

		// Kondisi untuk command
		switch command {
		case "/start":
			utils.StartCommand(chatID, bot)
			continue
		case "/close":
			hideMenu := telegrambot.NewRemoveKeyboard(true)
			msg := telegrambot.NewMessage(chatID, "Menu ditutup. Ketik /start untuk buka kembali.")
			msg.ReplyMarkup = hideMenu
			bot.Send(msg)
			continue
		case "/bantuan", "🆘Bantuan":
			utils.HelpCommand(chatID, bot)
			continue
		case "/harian", "📆Harian":
			utils.DailyTransactionCommand(chatID, bot, appCtx.TransactionService)
			continue;
		case "/bulanan", "📅Bulanan":
			utils.MonthlyTransactionCommand(chatID, bot, appCtx.TransactionService)
			continue;
		case "/hapus", "🔥Hapus":
			if (command == "/hapus") {
				msg := telegrambot.NewMessage(chatID, "Mohon maaf fitur sedang perbaikan 🙏")
				bot.Send(msg)
				continue;
			}
			utils.DeleteTransactionCommand(chatID, textMessage, bot, appCtx.TransactionService)
			continue
		case "/daftar", "📝Daftar":
			if (command == "/daftar") {
				utils.RegisterComand(textMessage, chatID, bot, appCtx.UserService)
			} else if (command == "📝Daftar") {
				appCtx.UserStateStore[chatID] = "awaiting_name"
				msg := telegrambot.NewMessage(chatID, "Silakan ketik nama lengkap kamu untuk daftar:")
				bot.Send(msg)
			}
			continue;
		case "/dashboard", "📊Dashboard":
			utils.DashboardCommand(chatID, bot, appCtx.UserService)
			continue;
		default:
			if strings.HasPrefix(textMessage, "/") {
				msg := telegrambot.NewMessage(chatID, "❓ Command tidak dikenali. Gunakan /help untuk melihat daftar command.")
				bot.Send(msg)
				continue;
			}
		}

		// Handler untuk listen input user
		if appCtx.UserStateStore[chatID] == "awaiting_name" {
			fullName := textMessage
		
			user, err := appCtx.UserService.RegisterUser(chatID, fullName)
			if err != nil {
				if (strings.Contains(err.Error(), "user sudah terdaftar")) {
					msg := telegrambot.NewMessage(chatID, "👤 User sudah terdaftar, silakan login.")
					bot.Send(msg)
					delete(appCtx.UserStateStore, chatID)
					continue
				}
			}
		
			// Hapus state
			delete(appCtx.UserStateStore, chatID)
		
			msg := telegrambot.NewMessage(chatID, fmt.Sprintf("✅ Selamat datang " + user.Name + ", mau lanjut lihat dashboard?"))
			msg.ReplyMarkup = telegrambot.InlineKeyboardMarkup{
				InlineKeyboard: [][]telegrambot.InlineKeyboardButton{
					{
						telegrambot.NewInlineKeyboardButtonURL("📊 Buka Dashboard",
							config.AppConfig.DashboardUrl + "?ref=" + utils.EncodeChatID(chatID)),
					},
				},
			}
			bot.Send(msg)
			continue
		}

		// Cek otoritas user
		countTransactionUser, _ := appCtx.TransactionService.CountTransactionsById(chatID)
		if (MaxLimitHit <= countTransactionUser) {
			msg := telegrambot.NewMessage(chatID, "Maaf anda sudah melebihi limit, hubungi admin di https://t.me/Hirumakun.")
			bot.Send(msg)
			continue;
		}

		// Klasifikasi prompt & buat transaksi
		// result, fullResponse, _ := classifyTransaction(update.Message.Text)
		result, fullResponse, _ := hitChatGpt(textMessage)

		transactionType, err := utils.ParseTransactionType(result.TransactionType)
		if err != nil {
			// msgTxt := ""
			msg := telegrambot.NewMessage(chatID, "Transaksi gagal dibuat!")
			bot.Send(msg)
		}

		// Simpan ke database
		loc, _ := time.LoadLocation("Asia/Jakarta")
		jakartaTime := time.Now().In(loc)

		tx := model.Transaction{
			ChatID: chatID,
			OriginalText: update.Message.Text,
			TransactionType: transactionType,
			Amount: result.Amount,
			Category: utils.Slugify(result.Category),
			TransactionDate: jakartaTime,
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
