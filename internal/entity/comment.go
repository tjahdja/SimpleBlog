package entity

import (
	"time"
)

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"`
	PostID    uint      `json:"post_id"`
	UserID    uint      `json:"user_id"`
}
