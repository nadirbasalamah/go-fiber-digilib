package auth

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Username       string         `json:"username" gorm:"unique"`
	Email          string         `json:"email" gorm:"unique"`
	Password       string         `json:"-"`
	Address        string         `json:"address" gorm:"type:text"`
	ProvinceID     uint           `json:"province_id"`
	CityID         uint           `json:"city_id"`
	DistrictID     uint           `json:"district_id"`
	ProfilePicture string         `json:"profile_picture"`
	Role           string         `json:"role" gorm:"type:role"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type RegisterRequest struct {
	Username   string `json:"username" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8,containsNumber,containsSpecialCharacter"`
	Address    string `json:"address" validate:"required"`
	ProvinceID uint   `json:"province_id" validate:"required"`
	CityID     uint   `json:"city_id" validate:"required"`
	DistrictID uint   `json:"district_id" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
