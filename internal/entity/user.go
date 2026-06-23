package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Username  string `json:"username" binding:"required,min=3,max=20"`
	Password  string `json:"password,omitempty" binding:"required,min=6"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
