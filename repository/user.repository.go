package repository

import (
	"finance-bot/model"

	"gorm.io/gorm"
)

type UserRepository interface {
    GetByChatId(chatId int64) (*model.User, error)
    GetTransactions(uuid string) (*[]model.Transaction, error)
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
    return &userRepository{db}
}

func (r *userRepository) GetByChatId(chatId int64) (*model.User, error) {
    var user model.User
    if err := r.db.Where("chat_id = ?", chatId).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) GetTransactions(uuid string) (*[]model.Transaction, error) {
    var user model.User
    if err := r.db.Preload("Transactions").Where("uuid = ?", uuid).First(&user).Error; err != nil {
        return nil, err
    }
    return &user.Transactions, nil
}

