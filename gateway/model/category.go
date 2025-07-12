package model

import "time"

type CategoryType string

const (
	CategoryTypeINCOME  CategoryType = "INCOME"
	CategoryTypeEXPENSE CategoryType = "EXPENSE"
)

type Category struct {
	ID        uint          `json:"id"`
	Name      string        `json:"name"`
	Type      CategoryType  `json:"type"`
	ParentID  *uint         `json:"parent_id,omitempty"`
	UserID    *uint         `json:"user_id,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
