package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"finance-bot/model"
	"finance-bot/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetDailyReport(db *gorm.DB, chatID int64) ([]model.Transaction, error) {
	today := time.Now().Truncate(24 * time.Hour) // 00:00:00 hari ini
	tomorrow := today.Add(24 * time.Hour)        // 00:00:00 besok

	var transactions []model.Transaction
	err := db.Where("chat_id = ? AND transaction_date >= ? AND transaction_date < ?", chatID, today, tomorrow).
		Order("transaction_date asc").
		Find(&transactions).Error

	return transactions, err
}

func GetMonthlyReport(db *gorm.DB, chatID int64) ([]model.Transaction, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startOfNextMonth := startOfMonth.AddDate(0, 1, 0)

	var transactions []model.Transaction
	err := db.Where("chat_id = ? AND transaction_date >= ? AND transaction_date < ?", chatID, startOfMonth, startOfNextMonth).
		Order("transaction_date asc").
		Find(&transactions).Error

	return transactions, err
}

func DeleteTransactionByID(db *gorm.DB, transactionID uint, chatID int64) error {
	result := db.Where("id = ? AND chat_id = ?", transactionID, chatID).Delete(&model.Transaction{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("transaction not found or does not belong to this user")
	}

	return nil
}


func CountTransactionsById(db *gorm.DB, chatID int64) (int64, error) {
	var count int64
	err := db.Model(&model.Transaction{}).Where("chat_id = ?", chatID).Count(&count).Error
	return count, err
}

func RegisterUser(db *gorm.DB, chatID int64, name string) (model.User, error) {
	var existing model.User
	if err := db.Where("chat_id = ?", chatID).First(&existing).Error; err == nil {
		return existing, errors.New("user sudah terdaftar")
	}

	user := model.User{
		UUID:    uuid.New().String(),
		ChatID:  chatID,
		Name:    name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&user).Error; err != nil {
		return user, err
	}

	return user, nil;
}

func CheckUser(db *gorm.DB, chatID int64) bool {
	var user model.User
	err := db.Where("chat_id = ?", chatID).First(&user).Error
	return err == nil
}

func FormatDailyReport(transactions []model.Transaction) string {
	var (
		report       strings.Builder
		totalOut     float64
		totalIn      float64
	)

	report.WriteString("📊 *Laporan Hari Ini*\n")

	for _, tx := range transactions {
		var transactionType string
		if tx.TransactionType == "EXPENSE" {
			transactionType = "🔴"
		} else if tx.TransactionType == "INCOME" {
			transactionType = "🟢"
		}

		formatAmount := utils.FormatRupiah(tx.Amount)

		line := fmt.Sprintf("#%d %s %s %s\n", tx.ID, transactionType, formatAmount, tx.OriginalText)
		report.WriteString(line)

		if tx.TransactionType == "EXPENSE" {
			totalOut += tx.Amount
		} else if tx.TransactionType == "INCOME" {
			totalIn += tx.Amount
		}
	}

	report.WriteString("\n")
	report.WriteString(fmt.Sprintf("🟢 Total Pemasukan: %s\n", utils.FormatRupiah(totalIn)))
	report.WriteString(fmt.Sprintf("🔴 Total Pengeluaran: %s\n", utils.FormatRupiah(totalOut)))

	return report.String()
}

func FormatMonthlyReport(transactions []model.Transaction) string {
	var (
		report       strings.Builder
		totalOut     float64
		totalIn      float64
	)

	report.WriteString("📆 *Laporan Bulan Ini*\n")

	for _, tx := range transactions {
		var transactionType string
		if tx.TransactionType == "EXPENSE" {
			transactionType = "🔴"
		} else if tx.TransactionType == "INCOME" {
			transactionType = "🟢"
		}

		formatAmount := utils.FormatRupiah(tx.Amount)
		formatDate := utils.FormatDate(tx.TransactionDate)

		line := fmt.Sprintf("#%d %s %s  %s %s\n", tx.ID, transactionType, formatDate, formatAmount, tx.OriginalText)
		report.WriteString(line)

		if tx.TransactionType == "EXPENSE" {
			totalOut += tx.Amount
		} else if tx.TransactionType == "INCOME" {
			totalIn += tx.Amount
		}
	}

	report.WriteString("\n")
	report.WriteString(fmt.Sprintf("🟢 Total Pemasukan: %s\n", utils.FormatRupiah(totalIn)))
	report.WriteString(fmt.Sprintf("🔴 Total Pengeluaran: %s\n", utils.FormatRupiah(totalOut)))

	return report.String()
}





