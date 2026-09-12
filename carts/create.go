package carts

import (
	"context"
	"errors"
	"go-fiber-digilib/db/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type create struct {
	repository *gorm.DB
}

func (c create) Create(ctx context.Context, cartReq *CartRequest) (Cart, error) {
	cart := models.Cart{
		BookID:   cartReq.BookID,
		UserID:   cartReq.UserID,
		Quantity: cartReq.Quantity,
	}
	record := new(Cart)

	err := c.repository.Transaction(func(tx *gorm.DB) error {
		book := new(models.Book)
		if err := tx.WithContext(ctx).First(book, "id = ?", cart.BookID).Error; err != nil {
			return err
		}

		if book.Stock < cart.Quantity {
			return errors.New("book out of stock")
		}

		result := c.repository.WithContext(ctx).Create(&cart)

		if err := result.Error; err != nil {
			return err
		}

		if err := result.WithContext(ctx).Preload(clause.Associations).Last(record).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return Cart{}, err
	}

	return *record, nil
}
