package carts

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type getbyuser struct {
	repository *gorm.DB
}

func (g getbyuser) GetByUser(ctx context.Context, userId uint) ([]Cart, error) {
	carts := []Cart{}

	if err := g.repository.
		WithContext(ctx).
		Preload(clause.Associations).
		Where("user_id = ?", userId).
		Find(&carts).
		Error; err != nil {
		return nil, err
	}

	return carts, nil
}
