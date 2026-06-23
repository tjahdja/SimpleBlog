package entity

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Content   string         `gorm:"type:text;not null" json:"content" binding:"required,min=3,max=1000"`
	PostID    uint           `gorm:"not null" json:"post_id" binding:"required"`
	UserID    uint           `gorm:"not null" json:"user_id" binding:"required"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
