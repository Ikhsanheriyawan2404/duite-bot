package repository

import (
    "errors"
	"time"
    
	"finance-bot/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
    GetByChatId(chatId int64) (*model.User, error)
    GetTransactions(uuid string) (*[]model.Transaction, error)
    RegisterUser(chatID int64, name string) (model.User, error)
    CheckUser(chatID int64) bool
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
    return &userRepository{db}
}

func (r *userRepository) GetByChatId(chatId int64) (*model.User, error) {
    var user model.User
    if err := r.db.
        Select("id", "name", "uuid", "chat_id").
        Where("chat_id = ?", chatId).
        First(&user).Error; err != nil {
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

func (r *userRepository) RegisterUser(chatID int64, name string) (model.User, error) {
	var existing model.User
	if err := r.db.Where("chat_id = ?", chatID).First(&existing).Error; err == nil {
		return existing, errors.New("user sudah terdaftar")
	}

	user := model.User{
		UUID:    uuid.New().String(),
		ChatID:  chatID,
		Name:    name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := r.db.Create(&user).Error; err != nil {
		return user, err
	}

	return user, nil;
}

func (r *userRepository) CheckUser(chatID int64) bool {
	var user model.User
	err := r.db.Where("chat_id = ?", chatID).First(&user).Error
	return err == nil
}
