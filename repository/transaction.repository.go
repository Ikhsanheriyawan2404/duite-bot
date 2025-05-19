package repository

import (
	"errors"
	"time"
    
	"finance-bot/model"
    
	"gorm.io/gorm"
)


type TransactionRepository interface {
    GetDailyReport(chatID int64) ([]model.Transaction, error)
    GetMonthlyReport(chatID int64) ([]model.Transaction, error)
    DeleteTransactionByID(transactionID uint, chatID int64) error
    CountTransactionsById(chatID int64) (int64, error)
}

type transactionRepository struct {
    db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *transactionRepository {
    return &transactionRepository{db}
}

func (r *transactionRepository) GetDailyReport(chatID int64) ([]model.Transaction, error) {
	today := time.Now().Truncate(24 * time.Hour) // 00:00:00 hari ini
	tomorrow := today.Add(24 * time.Hour)        // 00:00:00 besok

	var transactions []model.Transaction
	err := r.db.Where("chat_id = ? AND transaction_date >= ? AND transaction_date < ?", chatID, today, tomorrow).
		Order("transaction_date asc").
		Find(&transactions).Error

	return transactions, err
}

func (r *transactionRepository) GetMonthlyReport(chatID int64) ([]model.Transaction, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startOfNextMonth := startOfMonth.AddDate(0, 1, 0)

	var transactions []model.Transaction
	err := r.db.Where("chat_id = ? AND transaction_date >= ? AND transaction_date < ?", chatID, startOfMonth, startOfNextMonth).
		Order("transaction_date asc").
		Find(&transactions).Error

	return transactions, err
}

func (r *transactionRepository) DeleteTransactionByID(transactionID uint, chatID int64) error {
	result := r.db.Where("id = ? AND chat_id = ?", transactionID, chatID).Delete(&model.Transaction{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("transaction not found or does not belong to this user")
	}

	return nil
}

func (r *transactionRepository) CountTransactionsById(chatID int64) (int64, error) {
	var count int64
	err := r.db.Model(&model.Transaction{}).Where("chat_id = ?", chatID).Count(&count).Error
	return count, err
}