package repository

import (
	"finance-bot/model"
    
	"gorm.io/gorm"
)


type TransactionRepository interface {
    CountTransactionsById(chatID int64) (int64, error)
	CreateTransaction(tx *model.Transaction) error
}

type transactionRepository struct {
    db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *transactionRepository {
    return &transactionRepository{db}
}

func (r *transactionRepository) CountTransactionsById(chatID int64) (int64, error) {
	var count int64
	err := r.db.Model(&model.Transaction{}).Where("chat_id = ?", chatID).Count(&count).Error
	return count, err
}

func (r *transactionRepository) CreateTransaction(tx *model.Transaction) error {
	return r.db.Create(tx).Error
}