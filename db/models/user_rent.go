package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRent struct {
	ID        uint           `gorm:"primaryKey"`
	RentID    uint           `json:"rent_id"`
	Rent      Rent           `json:"rent"`
	CartID    uint           `json:"cart_id"`
	Cart      Cart           `json:"cart"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
