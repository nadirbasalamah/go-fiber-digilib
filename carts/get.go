package carts

import (
	"context"

	"gorm.io/gorm"
)

type get struct {
	repository *gorm.DB
}

func (g get) GetByID(ctx context.Context, id uint) (Cart, error) {
	cart := new(Cart)

	if err := g.repository.WithContext(ctx).First(cart, "id = ?", id).Error; err != nil {
		return Cart{}, err
	}

	return *cart, nil
}
