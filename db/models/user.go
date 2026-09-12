package models

import (
	"database/sql/driver"
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
	Role           Role           `json:"role" gorm:"type:role"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type Role string

const (
	Enduser Role = "user"
	Admin   Role = "admin"
)

func (p *Role) Scan(value any) error {
	*p = Role(value.([]byte))
	return nil
}

func (p Role) Value() (driver.Value, error) {
	return string(p), nil
}
