package model

import (
	"time"
)

type CategoryType string

const (
	CategoryTypeINCOME  CategoryType = "INCOME"
	CategoryTypeEXPENSE CategoryType = "EXPENSE"
)

type Category struct {
	ID        uint         `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name      string       `gorm:"type:varchar(100);not null" json:"name"`
	Type      CategoryType `gorm:"type:varchar(10);not null;check:type IN ('INCOME','EXPENSE')" json:"type"`
	ParentID  *uint        `gorm:"type:bigint" json:"parent_id"`
	UserID    *uint        `gorm:"type:bigint" json:"user_id"`

	Parent     *Category   `gorm:"foreignKey:ParentID"`
	User       *User       `gorm:"foreignKey:UserID"`
	Children   []Category  `gorm:"foreignKey:ParentID"`
	CreatedAt  time.Time   `gorm:"autoCreateTime;type:timestamptz;default:now()" json:"created_at"`
	UpdatedAt  time.Time   `gorm:"autoUpdateTime;type:timestamptz;default:now()" json:"updated_at"`
}
