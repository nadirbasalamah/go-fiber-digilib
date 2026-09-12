package users

import (
	"mime/multipart"
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

type EditProfileRequest struct {
	Username       string `form:"username" validate:"required"`
	Email          string `form:"email" validate:"required,email"`
	Address        string `form:"address" validate:"required"`
	ProfilePicture string
	File           *multipart.FileHeader
}
