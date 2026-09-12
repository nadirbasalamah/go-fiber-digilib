package rents

import (
	"context"
	"go-fiber-digilib/pkg/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type getall struct {
	repository *gorm.DB
}

func (g getall) GetAll(ctx context.Context, pagination utils.Pagination) (utils.Pagination, error) {
	rents := []UserRent{}

	if err := g.repository.
		WithContext(ctx).
		Joins("JOIN rents ON rents.id = user_rents.rent_id").
		Joins("JOIN users ON users.id = rents.user_id").
		Scopes(utils.Paginate(&rents, &pagination, g.repository)).
		Preload(clause.Associations).
		Preload("Cart." + clause.Associations).
		Preload("Cart.Book." + clause.Associations).
		Preload("Rent." + clause.Associations).
		Find(&rents).Error; err != nil {
		return utils.Pagination{}, err
	}

	pagination.Rows = rents

	return pagination, nil
}
