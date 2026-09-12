package carts

import (
	"go-fiber-digilib/books"
	"go-fiber-digilib/users"
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	BookID    uint           `json:"book_id"`
	Book      books.Book     `json:"book"`
	UserID    uint           `json:"user_id"`
	User      users.User     `json:"user"`
	Quantity  uint           `json:"quantity"`
	IsRented  bool           `json:"is_rented"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type CartRequest struct {
	BookID   uint `json:"book_id" validate:"required"`
	Quantity uint `json:"quantity" validate:"required,gte=1"`
	UserID   uint
}
