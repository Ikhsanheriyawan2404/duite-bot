package service

import (
	"finance-bot/model"
	"finance-bot/repository"
)

type UserService interface {
    GetUser(chatId int64) (*model.User, error)
    GetTransactions(uuid string) (*[]model.Transaction, error)
    RegisterUser(chatID int64, name string) (model.User, error)
    CheckUser(chatID int64) bool
}

type userService struct {
    userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
    return &userService{userRepo}
}

func (s *userService) GetUser(chatId int64) (*model.User, error) {
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