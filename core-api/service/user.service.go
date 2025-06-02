package service

import (
	"context"
	"finance-bot/model"
	"finance-bot/repository"
	"finance-bot/utils"
	"time"

	redis "finance-bot/config"
)

type UserService interface {
    GetByChatId(chatId int64) (*model.User, error)
    GetTransactions(uuid string) (*[]model.Transaction, error)
    RegisterUser(chatID int64, name string) (model.User, error)
    CheckUser(chatID int64) bool
	GetDailyReport(chatID int64) ([]model.Transaction, error)
    GetMonthlyReport(chatID int64) ([]model.Transaction, error)
    DeleteTransactionByID(transactionID uint, chatID int64) error
    GenerateMagicLoginToken(chatID int64) (string, error)
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

func (s *userService) GenerateMagicLoginToken(chatID int64) (string, error) {
	user, _ := s.userRepo.GetByChatId(chatID)

	token := utils.GenerateRandomToken(32)
	key := "magic_login:" + token

	redis.Client.Set(context.Background(), key, user.ChatID, 15*time.Minute).Err()

	return token, nil
}