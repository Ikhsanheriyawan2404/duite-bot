package service

import (
	"finance-bot/model"
	"finance-bot/repository"
)

type TransactionService interface {
    CountTransactionsById(chatID int64) (int64, error)
	CreateTransaction(tx *model.Transaction) error
	GetTransactionWithCategory(id uint) (*model.Transaction, error)
}

type transactionService struct {
    transactionRepo repository.TransactionRepository
}

func NewTransactionService(transactionRepo repository.TransactionRepository) TransactionService {
    return &transactionService{transactionRepo}
}

func (s *transactionService) CountTransactionsById(chatID int64) (int64, error) {
    return s.transactionRepo.CountTransactionsById(chatID)
}

func (s *transactionService) CreateTransaction(tx *model.Transaction) error {
    return s.transactionRepo.CreateTransaction(tx)
}

func (s *transactionService) GetTransactionWithCategory(id uint) (*model.Transaction, error) {
    return s.transactionRepo.GetTransactionWithCategory(id)
}





