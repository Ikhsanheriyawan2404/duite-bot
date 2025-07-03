
package model

import "time"

type Category struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	ParentID  *uint     `json:"parent_id"`
	UserID    *uint     `json:"user_id"`
	Parent    *Category `json:"parent,omitempty"`
	User      *User     `json:"user,omitempty"`
	Children  []Category `json:"children,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}