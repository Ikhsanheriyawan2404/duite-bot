package model

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	UUID      string    `gorm:"type:char(36);uniqueIndex;not null" json:"uuid"`
	ChatID    int64     `gorm:"type:bigint;uniqueIndex;not null" json:"chat_id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	IsPaid    *time.Time `gorm:"type:timestamp" json:"is_paid"`
	CreatedAt time.Time `gorm:"autoCreateTime;type:timestamptz;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;type:timestamptz;default:now()" json:"updated_at"`

	Categories   []Category   `gorm:"foreignKey:UserID"`
	Transactions []Transaction `gorm:"foreignKey:ChatID;references:ChatID"`
}
