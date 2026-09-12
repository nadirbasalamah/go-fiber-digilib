package carts

import (
	"context"
	"errors"
	"go-fiber-digilib/db/models"
	"go-fiber-digilib/pkg/utils"

	"gorm.io/gorm"
)

type update struct {
	repository *gorm.DB
	get
}

func (u update) Update(ctx context.Context, cartReq *CartRequest, id uint) (Cart, error) {
	cart := models.Cart{
		BookID:   cartReq.BookID,
		Quantity: cartReq.Quantity,
	}

	err := u.repository.Transaction(func(tx *gorm.DB) error {
		book := new(models.Book)
		if err := tx.WithContext(ctx).First(book, "id = ?", cart.BookID).Error; err != nil {
			return err
		}

		if book.Stock < cart.Quantity {
			return errors.New("book out of stock")
		}

		result := u.repository.WithContext(ctx).Scopes(utils.CurrentUser(cartReq.UserID)).Where("id = ?", id).Updates(&cart)

		isFailed := result.Error != nil || result.RowsAffected == 0

		if isFailed {
			return errors.New("update failed")
		}

		return nil
	})

	if err != nil {
		return Cart{}, err
	}

	record, err := u.get.GetByID(ctx, id)

	if err != nil {
		return Cart{}, err
	}

	return record, nil
}
