package entity

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"type:varchar(255);not null" json:"title" binding:"required,min=5,max=255"`
	Content   string         `gorm:"type:text;not null" json:"content" binding:"required,min=10"`
	AuthorID  uint           `gorm:"not null" json:"author_id" binding:"required"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // json:"-" hides it from API responses
}
