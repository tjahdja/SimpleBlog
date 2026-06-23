package entity

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username" binding:"required,min=3,max=20"`
	Password  string    `json:"password,omitempty" binding:"required,min=6"`
}
