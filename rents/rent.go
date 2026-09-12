package rents

import (
	"go-fiber-digilib/carts"
	"go-fiber-digilib/users"
	"time"

	"gorm.io/gorm"
)

type Rent struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id"`
	User       users.User     `json:"user"`
	Quantity   uint           `json:"quantity"`
	Fee        float64        `json:"fee"`
	Courier    string         `json:"courier"`
	Duration   uint           `json:"duration"`
	Status     string         `json:"rent_status" gorm:"type:rent_status"`
	ReturnTime time.Time      `json:"return_time"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	ReturnedAt time.Time      `json:"returned_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type UserRent struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	RentID    uint           `json:"rent_id"`
	Rent      Rent           `json:"rent"`
	CartID    uint           `json:"cart_id"`
	Cart      carts.Cart     `json:"cart"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type RentRequest struct {
	CartItems  []uint `json:"cart_items" validate:"required"`
	Duration   uint   `json:"duration" validate:"required,gte=1"`
	Courier    string `json:"courier" validate:"required,validCourier"`
	Fee        float64
	UserID     uint
	ReturnTime time.Time
}

type RentUpdateRequest struct {
	Status string `json:"status" validate:"required"`
}
