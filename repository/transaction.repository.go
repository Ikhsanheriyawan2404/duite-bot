package repository

import "gorm.io/gorm"


type TransactionRepository interface {
    
}

type transactionRepository struct {
    db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *transactionRepository {
    return &transactionRepository{db}
}