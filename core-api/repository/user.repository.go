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
	GetDailyReport(chatID int64) ([]model.Transaction, error)
    GetMonthlyReport(chatID int64) ([]model.Transaction, error)
    DeleteTransactionByID(transactionID uint, chatID int64) error
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

func (r *userRepository) GetDailyReport(chatID int64) ([]model.Transaction, error) {
	today := time.Now().Truncate(24 * time.Hour) // 00:00:00 hari ini
	tomorrow := today.Add(24 * time.Hour)        // 00:00:00 besok

	var transactions []model.Transaction
	err := r.db.
		Select("id", "chat_id", "amount", "category", "transaction_date", "original_text", "transaction_type").
		Where("chat_id = ? AND transaction_date >= ? AND transaction_date < ?", chatID, today, tomorrow).
		Order("transaction_date asc").
		Find(&transactions).Error

	return transactions, err
}

func (r *userRepository) GetMonthlyReport(chatID int64) ([]model.Transaction, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startOfNextMonth := startOfMonth.AddDate(0, 1, 0)

	var transactions []model.Transaction
	err := r.db.
		Select("id", "chat_id", "amount", "category", "transaction_date", "original_text", "transaction_type").
		Where("chat_id = ? AND transaction_date >= ? AND transaction_date < ?", chatID, startOfMonth, startOfNextMonth).
		Order("transaction_date asc").
		Find(&transactions).Error

	return transactions, err
}

func (r *userRepository) DeleteTransactionByID(transactionID uint, chatID int64) error {
	result := r.db.Where("id = ? AND chat_id = ?", transactionID, chatID).Delete(&model.Transaction{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("transaction not found or does not belong to this user")
	}

	return nil
}