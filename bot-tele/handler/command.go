// handler/command_handler.go
package handler

import (
	"bot-tele/config"
	"bot-tele/model"
	"bot-tele/service"
	"bot-tele/static"
	"bot-tele/utils"
	"strconv"

	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type LLMResult struct {
	Result struct {
		TransactionType string  `json:"type"`
		Amount          float64 `json:"amount"`
		Category        string  `json:"category"`
		Date            string  `json:"date"`
	} `json:"result"`
}

type CommandHandler struct {
	APIClient *service.APIClient
}

func HandleCommandAndInput(update tgbotapi.Update, bot *tgbotapi.BotAPI, apiClient *service.APIClient) {
	chatID := update.Message.Chat.ID
	inputMessage := update.Message.Text
	parts := strings.Fields(inputMessage)
	if len(parts) == 0 {
		return
	}
	command := parts[0]

	switch command {
	case "/start":
		handleStart(chatID, bot)
	case "/close":
		handleCloseMenu(chatID, bot)
	case "/bantuan", "🆘Bantuan":
		handleHelp(chatID, bot)
	case "/harian", "📆Harian":
		handleDailyReport(chatID, bot)
	case "/bulanan", "📅Bulanan":
		handleMonthlyReport(chatID, bot)
	case "/hapus", "🔥Hapus":
		if (command == "🔥Hapus") {
			msg := tgbotapi.NewMessage(chatID, "Mohon maaf fitur sedang perbaikan 🙏")
			bot.Send(msg)
			return
		}
		handleDeleteTransaction(chatID, inputMessage, bot)
		return
	case "/daftar", "📝Daftar":
		if (command == "📝Daftar") {
			msg := tgbotapi.NewMessage(chatID, "Mohon maaf fitur sedang perbaikan 🙏")
			bot.Send(msg)
			return
		}
		handleRegister(chatID, inputMessage, bot)
	case "/dashboard", "📊Dashboard":
		handleDashboard(chatID, bot)
	default:
		if strings.HasPrefix(inputMessage, "/") {
			msg := tgbotapi.NewMessage(chatID, static.WrongCommandText)
			bot.Send(msg)
			return
		}
		handleTransactionInput(chatID, inputMessage, bot)
	}
}

func handleStart(chatID int64, bot *tgbotapi.BotAPI) {
	menu := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📝Daftar"),
			tgbotapi.NewKeyboardButton("📊Dashboard"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📆Harian"),
			tgbotapi.NewKeyboardButton("📅Bulanan"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🆘Bantuan"),
			tgbotapi.NewKeyboardButton("🔥Hapus"),
		),
	)
	msg := tgbotapi.NewMessage(chatID, static.WelcomeText)
	msg.ReplyMarkup = menu
	bot.Send(msg)
}

func handleCloseMenu(chatID int64, bot *tgbotapi.BotAPI) {
	hideMenu := tgbotapi.NewRemoveKeyboard(true)
	msg := tgbotapi.NewMessage(chatID, static.CloseMenuText)
	msg.ReplyMarkup = hideMenu
	bot.Send(msg)
}

func handleHelp(chatID int64, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(chatID, static.HelpText)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func handleDailyReport(chatID int64, bot *tgbotapi.BotAPI) {
	var transactions []model.Transaction
	endpoint := fmt.Sprintf("/users/%d/transactions/daily", chatID)
	err := service.Client.Request("GET", endpoint, nil, &transactions)
	if err != nil {
		log.Println("Error:", err)
	}
	report := utils.FormatDailyReport(transactions)
	msg := tgbotapi.NewMessage(chatID, utils.EscapeMarkdown(report))
	msg.ParseMode = "MarkdownV2"
	bot.Send(msg)
}

func handleMonthlyReport(chatID int64, bot *tgbotapi.BotAPI) {
	var transactions []model.Transaction
	endpoint := fmt.Sprintf("/users/%d/transactions/monthly", chatID)
	err := service.Client.Request("GET", endpoint, nil, &transactions)
	if err != nil {
		log.Println("Error:", err)
	}
	report := utils.FormatMonthlyReport(transactions)
	msg := tgbotapi.NewMessage(chatID, utils.EscapeMarkdown(report))
	msg.ParseMode = "MarkdownV2"
	bot.Send(msg)
}

func handleDeleteTransaction(chatID int64, inputMessage string, bot *tgbotapi.BotAPI) {
	type Response struct {
		Message string `json:"message"`
	}
	var response Response

	parts := strings.Fields(inputMessage)
	if len(parts) < 2 {
		msg := tgbotapi.NewMessage(chatID, "Format salah. Gunakan: /hapus {ID transaksi}")
		bot.Send(msg)
		return
	}
	
	indexStr := parts[1]
	transactionId, err := strconv.Atoi(indexStr)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "ID transaksi harus berupa angka.")
		bot.Send(msg)
		return
	}

	endpoint := fmt.Sprintf("/users/%d/transactions/%d", chatID, transactionId)
	err = service.Client.Request("DELETE", endpoint, nil, &response)
	msg := tgbotapi.NewMessage(chatID, response.Message)
	bot.Send(msg)
}

func handleRegister(chatID int64, inputMessage string, bot *tgbotapi.BotAPI) {
	parts := strings.Fields(inputMessage)
	if len(parts) < 2 {
		msg := tgbotapi.NewMessage(chatID, "Format salah. Gunakan: /daftar {nama}")
		bot.Send(msg)
		return
	}

	userFullName := strings.Join(parts[1:], " ")
	reqBody := map[string]any{
		"chat_id": chatID,
		"name":    userFullName,
	}
	var user model.User
	err := service.Client.Request("POST", "/users/register", reqBody, &user)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Anda sudah terdaftar")
		bot.Send(msg)
		return
	}
	msg := tgbotapi.NewMessage(chatID, "✅ Selamat datang " + user.Name + ", mau lanjut lihat dashboard?")
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{
				tgbotapi.NewInlineKeyboardButtonURL("📊 Buka Dashboard",
					config.AppConfig.DashboardUrl + "?ref=" + utils.EncodeChatID(chatID)),
			},
		},
	}
	bot.Send(msg)
}

func handleDashboard(chatID int64, bot *tgbotapi.BotAPI) {
	type UserExistResponse struct {
		Exist bool `json:"exist"`
	}

	var res UserExistResponse
	url := fmt.Sprintf("/users/%d/exists", chatID)
	service.Client.Request("GET", url, nil, &res)
	if !res.Exist {
		msg := tgbotapi.NewMessage(chatID, "Yuk, daftar dulu biar bisa lanjut")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(chatID, "Klik tombol di bawah untuk membuka dashboard:")
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{
				tgbotapi.NewInlineKeyboardButtonURL("📊 Buka Dashboard",
					config.AppConfig.DashboardUrl + "?ref=" + utils.EncodeChatID(chatID)),
			},
		},
	}
	bot.Send(msg)
}

func handleTransactionInput(chatID int64, inputMessage string, bot *tgbotapi.BotAPI) {
	// 1. Kirim permintaan ke LLM API
	fullPrompt := fmt.Sprintf(static.PromptDefault, inputMessage)
	var llmResp LLMResult
	err := service.Client.Request(
		"POST",
		fmt.Sprintf("/users/%d/transactions/ai-classify", chatID),
		map[string]string{"prompt": fullPrompt}, // request body
		&llmResp,
	)

	if err != nil {
		log.Println("Gagal memanggil API klasifikasi:", err)
		bot.Send(tgbotapi.NewMessage(chatID, "Terjadi kesalahan saat memproses pesanmu."))
		return
	}

	// 2. Format hasil klasifikasi
	// reply := fmt.Sprintf(
	// 	"✅ Transaksi berhasil d:\n\n📄 *Jenis:* %s\n💰 *Jumlah:* %.2f\n🏷️ *Kategori:* %s\n📅 *Tanggal:* %s",
	// 	llmResp.Result.TransactionType,
	// 	llmResp.Result.Amount,
	// 	llmResp.Result.Category,
	// 	llmResp.Result.Date,
	// )

	// 2. Format hasil klasifikasi
	reply := fmt.Sprintf(
		"✅ Transaksi berhasil dicatat!\n\n📂 Tipe: %s\n💰 Jumlah: %s\n🏷️ Kategori: %s\n",
		llmResp.Result.TransactionType,
		llmResp.Result.Amount,
		llmResp.Result.Category,
	)

	// 3. Kirim balasan ke Telegram
	msg := tgbotapi.NewMessage(chatID, reply)
	bot.Send(msg)
}

