// handler/command_handler.go
package handler

import (
	"bot-tele/config"
	"bot-tele/model"
	"bot-tele/service"
	"bot-tele/static"
	"bot-tele/utils"

	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var appCtx = &model.AppContext{
    UserStateStore: make(map[int64]string),
}

func HandleCommandAndInput(update tgbotapi.Update, bot *tgbotapi.BotAPI, apiClient *service.APIClient) {
	chatID := update.Message.Chat.ID
	inputMessage := update.Message.Text
	parts := strings.Fields(inputMessage)
	if len(parts) == 0 {
		return
	}
	command := parts[0]
	
	if !utils.IsCommand(inputMessage) {
		state := appCtx.UserStateStore[chatID]
		if state != "" {
			handleUserStateInput(chatID, inputMessage, bot, apiClient)
			return
		}
	} else {
		delete(appCtx.UserStateStore, chatID)
	}

	switch command {
	case "/start":
		handleStart(chatID, bot)
	case "/close":
		handleCloseMenu(chatID, bot)
	case "/bantuan", "🆘Bantuan":
		handleHelp(chatID, bot)
	case "/harian", "📆Harian":
		handleDailyReport(chatID, bot, apiClient)
	case "/bulanan", "📅Bulanan":
		handleMonthlyReport(chatID, bot, apiClient)
	case "/hapus", "🔥Hapus":
		if (command == "🔥Hapus") {
			appCtx.UserStateStore[chatID] = "awaiting_transaction_id"
			msg := tgbotapi.NewMessage(chatID, "💥 Kirim ID transaksi yang mau kamu hapus, biar aku bantu hapuskan.")
			bot.Send(msg)
			return
		}

		parts := strings.Fields(inputMessage)
		if len(parts) < 2 {
			bot.Send(tgbotapi.NewMessage(chatID, "⚠️ Format salah nih! Coba ketik kayak gini:\n`/hapus 123`"))
			return
		}

		inputMessage = parts[1]

		handleDeleteTransaction(chatID, inputMessage, bot, apiClient)
		return
	case "/daftar", "📝Daftar":
		if (command == "📝Daftar") {
			appCtx.UserStateStore[chatID] = "awaiting_register_name"
			msg := tgbotapi.NewMessage(chatID, "📝 Ketik nama lengkap kamu ya, biar aku bisa daftarin.")
			bot.Send(msg)
			return
		}

		if len(parts) < 2 {
			msg := tgbotapi.NewMessage(chatID, "Waduuhh, kamu belum isi nama nih 😅\nCoba ketik kayak gini ya:\n👉 /daftar Udin Andria")
			bot.Send(msg)
			return
		}

		inputMessage = strings.Join(parts[1:], " ")

		handleRegister(chatID, inputMessage, bot, apiClient)
	case "/dashboard", "📊Dashboard":
		handleDashboard(chatID, bot, apiClient)
	default:
		if strings.HasPrefix(inputMessage, "/") {
			msg := tgbotapi.NewMessage(chatID, static.WrongCommandText)
			bot.Send(msg)
			return
		}
		handleTransactionInput(chatID, inputMessage, bot, apiClient)
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

func handleDailyReport(chatID int64, bot *tgbotapi.BotAPI, apiClient *service.APIClient) {
	var transactions []model.Transaction
	endpoint := fmt.Sprintf("/users/%d/transactions/daily", chatID)
	err := apiClient.Request("GET", endpoint, nil, &transactions)
	if err != nil {
		log.Println("Error:", err)
	}
	report := utils.FormatDailyReport(transactions)
	msg := tgbotapi.NewMessage(chatID, utils.EscapeMarkdown(report))
	msg.ParseMode = "MarkdownV2"
	bot.Send(msg)
}

func handleMonthlyReport(chatID int64, bot *tgbotapi.BotAPI, apiClient *service.APIClient) {
	var transactions []model.Transaction
	endpoint := fmt.Sprintf("/users/%d/transactions/monthly", chatID)
	err := apiClient.Request("GET", endpoint, nil, &transactions)
	if err != nil {
		log.Println("Error:", err)
	}
	report := utils.FormatMonthlyReport(transactions)
	msg := tgbotapi.NewMessage(chatID, utils.EscapeMarkdown(report))
	msg.ParseMode = "MarkdownV2"
	bot.Send(msg)
}

func handleDeleteTransaction(chatID int64, inputMessage string, bot *tgbotapi.BotAPI, apiClient *service.APIClient) {
	type Response struct {
		Message string `json:"message"`
	}
	var response Response
	
	indexStr := inputMessage
	transactionId, err := strconv.Atoi(indexStr)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "🧠 ID transaksi harus berupa angka yaa."))
		return
	}

	endpoint := fmt.Sprintf("/users/%d/transactions/%d", chatID, transactionId)
	apiClient.Request("DELETE", endpoint, nil, &response)
	msg := tgbotapi.NewMessage(chatID, response.Message)
	bot.Send(msg)
}

func handleRegister(chatID int64, inputMessage string, bot *tgbotapi.BotAPI, apiClient *service.APIClient) {
	userFullName := inputMessage
	reqBody := map[string]any{
		"chat_id": chatID,
		"name":    userFullName,
	}
	var user model.User
	err := apiClient.Request("POST", "/users/register", reqBody, &user)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Eh, btw kamu udah daftar sebelumnya, hehe")
		msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
				{
					tgbotapi.NewInlineKeyboardButtonURL("📊 Mau aku bantu lihat dashboard?",
						config.AppConfig.DashboardUrl + "?ref=" + utils.EncodeChatID(chatID)),
				},
			},
		}
		bot.Send(msg)
		return
	}
	msg := tgbotapi.NewMessage(chatID, "Hai " + user.Name + ", mau aku bantu lihat dashboard?")
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

func handleDashboard(chatID int64, bot *tgbotapi.BotAPI, apiClient *service.APIClient) {
	type UserExistResponse struct {
		Exist bool `json:"exist"`
	}

	var res UserExistResponse
	url := fmt.Sprintf("/users/%d/exists", chatID)
	apiClient.Request("GET", url, nil, &res)
	if !res.Exist {
		msg := tgbotapi.NewMessage(chatID, "Yuk, daftar dulu biar bisa lanjut!\nGampang kok, tinggal ketik: daftar NamaKamu\nContoh: daftar Budi")

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

func handleTransactionInput(chatID int64, inputMessage string, bot *tgbotapi.BotAPI, apiClient *service.APIClient) {
	// Kirim permintaan ke API klasifikasi & penyimpanan transaksi sekaligus
	reqBody := map[string]string{
		"prompt": inputMessage,
	}

	var apiResponse struct {
		Message string           `json:"message"`
		Usage   any              `json:"usage"`
		Data    model.Transaction `json:"data"`
		Error   string             `json:"error,omitempty"` // opsional, bisa kosong
	}

	err := apiClient.Request(
		"POST",
		fmt.Sprintf("/users/%d/transactions/ai-classify", chatID),
		reqBody,
		&apiResponse,
	)

	

	if err != nil {
		if apiResponse.Error != "" {
			bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ %s", apiResponse.Error)))
			return
		}
		// Balas error dari API langsung ke Telegram
		log.Println("❌ Gagal memproses transaksi:", err)
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("⚠️ %s", err.Error())))
		return
	}

	// Format dan kirim balasan ke Telegram
	reply := fmt.Sprintf(
		"✅ Transaksi berhasil dicatat!\n\n📂 Tipe: %s\n💰 Jumlah: %s\n🏷️ Kategori: %s\n🗓️ Tanggal: %s",
		apiResponse.Data.TransactionType,
		utils.FormatRupiah(apiResponse.Data.Amount),
		apiResponse.Data.Category,
		apiResponse.Data.TransactionDate.Format("2006-01-02"),
	)

	msg := tgbotapi.NewMessage(chatID, reply)
	bot.Send(msg)
}

func handleUserStateInput(chatID int64, inputMessage string, bot *tgbotapi.BotAPI, apiClient *service.APIClient) {
	state := appCtx.UserStateStore[chatID]

	switch state {
	case "awaiting_register_name":
		handleRegister(chatID, inputMessage, bot, apiClient)
		delete(appCtx.UserStateStore, chatID)
		return
	case "awaiting_transaction_id":
		handleDeleteTransaction(chatID, inputMessage, bot, apiClient)
		delete(appCtx.UserStateStore, chatID)
		return
	}
}



