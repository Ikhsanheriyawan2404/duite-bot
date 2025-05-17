package service

import (
	"finance-bot/model"
	"finance-bot/repository"
)

type UserService interface {
    GetUser(chatId int64) (*model.User, error)
    GetTransactions(uuid string) (*[]model.Transaction, error)
}

type userService struct {
    userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
    return &userService{userRepo}
}

func (u *userService) GetUser(chatId int64) (*model.User, error) {
    return u.userRepo.GetByChatId(chatId)
}

func (u *userService) GetTransactions(uuid string) (*[]model.Transaction, error) {
    return u.userRepo.GetTransactions(uuid)
}
