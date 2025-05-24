package service

import (
	"finance-bot/model"
	"finance-bot/repository"
)

type UserService interface {
    GetByChatId(chatId int64) (*model.User, error)
    GetTransactions(uuid string) (*[]model.Transaction, error)
    RegisterUser(chatID int64, name string) (model.User, error)
    CheckUser(chatID int64) bool
	GetDailyReport(chatID int64) ([]model.Transaction, error)
    GetMonthlyReport(chatID int64) ([]model.Transaction, error)
    DeleteTransactionByID(transactionID uint, chatID int64) error
}

type userService struct {
    userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
    return &userService{userRepo}
}

func (s *userService) GetByChatId(chatId int64) (*model.User, error) {
    return s.userRepo.GetByChatId(chatId)
}

func (s *userService) GetTransactions(uuid string) (*[]model.Transaction, error) {
    return s.userRepo.GetTransactions(uuid)
}

func (s *userService) RegisterUser(chatID int64, name string) (model.User, error) {
    return s.userRepo.RegisterUser(chatID, name)
}

func (s *userService) CheckUser(chatID int64) bool {
    return s.userRepo.CheckUser(chatID)
}

func (s *userService) GetDailyReport(chatID int64) ([]model.Transaction, error) {
    return s.userRepo.GetDailyReport(chatID)
}

func (s *userService) GetMonthlyReport(chatID int64) ([]model.Transaction, error) {
    return s.userRepo.GetMonthlyReport(chatID)
}

func (s *userService) DeleteTransactionByID(transactionID uint, chatID int64) error {
    return s.userRepo.DeleteTransactionByID(transactionID, chatID)
}