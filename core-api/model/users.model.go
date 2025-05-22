package model

import (
	"time"
)

type User struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID      string     `gorm:"type:char(36);uniqueIndex;not null" json:"uuid"`
	ChatID    int64      `gorm:"uniqueIndex;not null" json:"chat_id"`
	Name      string     `gorm:"type:varchar(100);not null" json:"name"`
	IsPaid    *time.Time `gorm:"type:timestamp" json:"is_paid"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	Transactions []Transaction `gorm:"foreignKey:ChatID;references:ChatID"`
}
