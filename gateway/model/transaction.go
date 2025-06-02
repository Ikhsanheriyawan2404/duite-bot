package model

import "time"

type Transaction struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	TransactionType string    `json:"transaction_type"`
	Amount          float64   `json:"amount"`
	Category        string    `json:"category"`
	Description     string    `json:"description"`
	TransactionDate time.Time `json:"transaction_date"`
	ChatID          string    `json:"chat_id"`
	OriginalText    string    `json:"original_text"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
