package rents

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type getbyuser struct {
	repository *gorm.DB
}

func (g getbyuser) GetByUser(ctx context.Context, userId uint) ([]UserRent, error) {
	rents := []UserRent{}

	if err := g.repository.
		WithContext(ctx).
		Joins("JOIN rents ON rents.id = user_rents.rent_id").
		Where("rents.user_id = ?", userId).
		Preload(clause.Associations).
		Preload("Cart." + clause.Associations).
		Preload("Cart.Book." + clause.Associations).
		Preload("Rent." + clause.Associations).
		Find(&rents).Error; err != nil {
		return nil, err
	}

	return rents, nil
}
