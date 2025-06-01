package model

import "time"

type User struct {
	ID        uint      `json:"id"`
	UUID      string    `json:"uuid"`
	ChatID    int64     `json:"chat_id"`
	Name      string    `json:"name"`
	IsPaid    *time.Time `json:"is_paid,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
