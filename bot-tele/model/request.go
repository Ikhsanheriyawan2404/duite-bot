package model

import (
	"time"
)

type TransactionDto struct {
	ChatID          int64     `json:"chat_id"`
	OriginalText    string    `json:"original_text"`
	TransactionType TransactionType `json:"transaction_type"`
	Amount          float64   `json:"amount"`
	Category        string    `json:"category"`
	TransactionDate time.Time `json:"transaction_date"`
}
