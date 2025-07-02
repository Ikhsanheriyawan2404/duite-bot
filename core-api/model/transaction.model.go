package model

import (
	"time"
)

type TransactionType string

const (
	TransactionTypeINCOME  TransactionType = "INCOME"
	TransactionTypeEXPENSE TransactionType = "EXPENSE"
)

type Transaction struct {
	ID              uint             `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	TransactionType TransactionType  `gorm:"type:varchar(20);not null;check:transaction_type IN ('INCOME','EXPENSE')" json:"transaction_type"`
	Amount          float64          `gorm:"type:numeric(15,2);not null" json:"amount"`
	CategoryID      *uint            `gorm:"type:bigint" json:"category_id"`
	Category        *Category        `gorm:"foreignKey:CategoryID"`
	Description     string           `gorm:"type:text" json:"description"`
	TransactionDate time.Time        `gorm:"type:timestamptz;not null" json:"transaction_date"`
	ChatID          int64            `gorm:"type:bigint" json:"chat_id"`
	OriginalText    string           `gorm:"type:text" json:"original_text"`
	CreatedAt       time.Time        `gorm:"autoCreateTime;type:timestamptz;default:now()" json:"created_at"`
	UpdatedAt       time.Time        `gorm:"autoUpdateTime;type:timestamptz;default:now()" json:"updated_at"`
}
