package model

import "time"

type TransactionType string

const (
	INCOME  TransactionType = "INCOME"
	EXPENSE TransactionType = "EXPENSE"
)

type Transaction struct {
	ID              uint            `json:"id"`
	TransactionType TransactionType `json:"transaction_type"`
	Amount          float64         `json:"amount"`
	CategoryID      *uint           `json:"category_id"`
	Category        *Category       `json:"category,omitempty"`
	Description     string          `json:"description"`
	TransactionDate time.Time       `json:"transaction_date"`
	ChatID          int64           `json:"chat_id"`
	OriginalText    string          `json:"original_text"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}