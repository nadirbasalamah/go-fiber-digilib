package models

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	ID        uint           `gorm:"primaryKey"`
	BookID    uint           `json:"book_id"`
	Book      Book           `json:"book"`
	UserID    uint           `json:"user_id"`
	User      User           `json:"user"`
	Quantity  uint           `json:"quantity"`
	IsRented  bool           `json:"is_rented"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
