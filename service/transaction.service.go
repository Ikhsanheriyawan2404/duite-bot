package service

import (
	"finance-bot/model"
	"finance-bot/repository"
)

type TransactionService interface {
	GetDailyReport(chatID int64) ([]model.Transaction, error)
    GetMonthlyReport(chatID int64) ([]model.Transaction, error)
    DeleteTransactionByID(transactionID uint, chatID int64) error
    CountTransactionsById(chatID int64) (int64, error)
}

type transactionService struct {
    transactionRepo repository.TransactionRepository
}

func NewTransactionService(transactionRepo repository.TransactionRepository) TransactionService {
    return &transactionService{transactionRepo}
}

func (s *transactionService) GetDailyReport(chatID int64) ([]model.Transaction, error) {
    return s.transactionRepo.GetDailyReport(chatID)
}

func (s *transactionService) GetMonthlyReport(chatID int64) ([]model.Transaction, error) {
    return s.transactionRepo.GetMonthlyReport(chatID)
}

func (s *transactionService) DeleteTransactionByID(transactionID uint, chatID int64) error {
    return s.transactionRepo.DeleteTransactionByID(transactionID, chatID)
}

func (s *transactionService) CountTransactionsById(chatID int64) (int64, error) {
    return s.transactionRepo.CountTransactionsById(chatID)
}





