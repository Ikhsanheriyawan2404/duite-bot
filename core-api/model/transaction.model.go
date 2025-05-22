package model

import (
	"time"
)

type TransactionType string

const (
	INCOME  TransactionType = "INCOME"
	EXPENSE TransactionType = "EXPENSE"
)

type Transaction struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TransactionType TransactionType `gorm:"type:varchar(20);not null;check (transaction_type IN ('INCOME', 'EXPENSE'))" json:"transaction_type"`
	Amount          float64   `gorm:"type:numeric(15,2);not null" json:"amount"`
	Category        string    `gorm:"type:varchar(100);not null" json:"category"`
	Description     string    `gorm:"type:text" json:"description"`
	TransactionDate time.Time `gorm:"type:timestamptz;not null" json:"transaction_date"`
	ChatID          int64     `gorm:"not null" json:"chat_id"`
	OriginalText    string    `gorm:"type:text" json:"original_text"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
