package utils

import (
	"finance-bot/config"
	"finance-bot/service"
	"strconv"
	"strings"

	telegrambot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartCommand(chatID int64, bot *telegrambot.BotAPI) {
	menu := telegrambot.NewReplyKeyboard(
		telegrambot.NewKeyboardButtonRow(
			telegrambot.NewKeyboardButton("📝Daftar"),
			telegrambot.NewKeyboardButton("📊Dashboard"),
		),
		telegrambot.NewKeyboardButtonRow(
			telegrambot.NewKeyboardButton("📆Harian"),
			telegrambot.NewKeyboardButton("📅Bulanan"),
		),
		telegrambot.NewKeyboardButtonRow(
			telegrambot.NewKeyboardButton("🆘Bantuan"),
			telegrambot.NewKeyboardButton("🔥Hapus"),
		),
	)
	msg := telegrambot.NewMessage(chatID, WelcomeText)
	msg.ReplyMarkup = menu
	bot.Send(msg)
}

func HelpCommand(chatID int64, bot *telegrambot.BotAPI) {
	msg := telegrambot.NewMessage(chatID, HelpText)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func DailyTransactionCommand(chatID int64, bot *telegrambot.BotAPI, transactionService service.TransactionService) {
	transactions, _ := transactionService.GetDailyReport(chatID)
	report := FormatDailyReport(transactions)
	msg := telegrambot.NewMessage(chatID, EscapeMarkdown(report))
	msg.ParseMode = "MarkdownV2"
	bot.Send(msg)
}

func MonthlyTransactionCommand(chatID int64, bot *telegrambot.BotAPI, transactionService service.TransactionService) {
	transactions, _ := transactionService.GetMonthlyReport(chatID)
	report := FormatMonthlyReport(transactions)
	msg := telegrambot.NewMessage(chatID, EscapeMarkdown(report))
	msg.ParseMode = "MarkdownV2"
	bot.Send(msg)
}

func RegisterComand(textMessage string, chatID int64, bot *telegrambot.BotAPI, useService service.UserService) {
	parts := strings.Fields(textMessage)
	if len(parts) < 2 {
		msg := telegrambot.NewMessage(chatID, "Format salah. Gunakan: /daftar {nama}")
		bot.Send(msg)
		return
	}

	userFullName := parts[1]
	user, _ := useService.RegisterUser(chatID, userFullName)
	msg := telegrambot.NewMessage(chatID, "✅ Selamat datang " + user.Name + ", mau lanjut lihat dashboard?")
	bot.Send(msg)
}

func RegisterMenu(textMessage string, chatID int64, bot *telegrambot.BotAPI, useService service.UserService) {
	parts := strings.Fields(textMessage)
	if len(parts) < 2 {
		msg := telegrambot.NewMessage(chatID, "Format salah. Gunakan: /daftar {nama}")
		bot.Send(msg)
		return
	}

	userFullName := parts[1]
	user, _ := useService.RegisterUser(chatID, userFullName)
	msg := telegrambot.NewMessage(chatID, "Selamat datang " + user.Name + ", lanjut lihat dashboard?")
	msg.ReplyMarkup = telegrambot.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegrambot.InlineKeyboardButton{
			{
				telegrambot.NewInlineKeyboardButtonURL("📊 Buka Dashboard",
					config.AppConfig.DashboardUrl + "?ref=" + EncodeChatID(chatID)),
			},
		},
	}
	bot.Send(msg)
}

func DashboardCommand(chatID int64, bot *telegrambot.BotAPI, useService service.UserService) {
	userExist := useService.CheckUser(chatID)
	if (!userExist) {
		msg := telegrambot.NewMessage(chatID, "Mohon melakukan daftar terlebih dahulu")
		bot.Send(msg)
		return
	}
	msg := telegrambot.NewMessage(chatID, "Klik tombol di bawah untuk membuka dashboard:")
	msg.ReplyMarkup = telegrambot.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegrambot.InlineKeyboardButton{
			{
				telegrambot.NewInlineKeyboardButtonURL("📊 Buka Dashboard",
					config.AppConfig.DashboardUrl + "?ref=" + EncodeChatID(chatID)),
			},
		},
	}
	bot.Send(msg)
}

func DeleteTransactionCommand(chatID int64, textMessage string, bot *telegrambot.BotAPI, transactionService service.TransactionService) {
	parts := strings.Fields(textMessage)
	if len(parts) < 2 {
		msg := telegrambot.NewMessage(chatID, "Format salah. Gunakan: /hapus {ID transaksi}")
		bot.Send(msg)
		return
	}
	
	indexStr := parts[1]
	transactionId, err := strconv.Atoi(indexStr)
	if err != nil {
		msg := telegrambot.NewMessage(chatID, "ID transaksi harus berupa angka.")
		bot.Send(msg)
		return
	}

	err = transactionService.DeleteTransactionByID(uint(transactionId), chatID)
	if (err != nil) {
		msg := telegrambot.NewMessage(chatID, "Transaksi tidak ditemukan")
		bot.Send(msg)
		return
	}
	msg := telegrambot.NewMessage(chatID, "Transaksi berhasil dihapus!")
	bot.Send(msg)
}