package rents

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

func (c create) Create(ctx context.Context, rentReq *RentRequest) (Rent, error) {
	rent := models.Rent{
		UserID:     rentReq.UserID,
		Duration:   rentReq.Duration,
		ReturnTime: rentReq.ReturnTime,
		Courier:    rentReq.Courier,
		Status:     models.Pending,
	}

	rentRecord := new(Rent)

	err := c.repository.Transaction(func(tx *gorm.DB) error {
		// find cart items
		carts := []models.Cart{}

		if err := tx.WithContext(ctx).
			Where("is_rented = false").
			Find(&carts, rentReq.CartItems).
			Error; err != nil {
			return err
		}

		if len(carts) == 0 {
			return errors.New("carts not found")
		}

		// calculate total quantity
		totalQty := calculateTotalQty(carts)

		// calculate final fee = quantity * rentReq.fee (update it the "rent")
		rent.Quantity = totalQty
		rent.Fee = float64(rent.Quantity) * rentReq.Fee

		// create rent record
		if err := tx.WithContext(ctx).
			Create(&rent).
			Error; err != nil {
			return err
		}

		if err := tx.WithContext(ctx).
			Preload(clause.Associations).
			Last(&rentRecord).
			Error; err != nil {
			return err
		}

		// insert user_rent record (from carts)
		userRents := generateUserRents(rentRecord.ID, carts)

		if err := tx.WithContext(ctx).
			Create(userRents).
			Error; err != nil {
			return err
		}

		// update book stock & cart status to is rented = true
		for _, cart := range carts {
			if err := updateBookStock(ctx, tx, cart.BookID, cart.Quantity); err != nil {
				return err
			}

			cart.IsRented = true

			if err := tx.WithContext(ctx).Where("id = ? ", cart.ID).Updates(&cart).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return Rent{}, err
	}

	return *rentRecord, nil
}

func calculateTotalQty(carts []models.Cart) uint {
	var total uint

	for _, c := range carts {
		total += c.Quantity
	}

	return total
}

func generateUserRents(rentID uint, carts []models.Cart) []*models.UserRent {
	records := make([]*models.UserRent, len(carts))

	for idx, c := range carts {
		records[idx] = &models.UserRent{
			RentID: rentID,
			CartID: c.ID,
		}
	}

	return records
}

func updateBookStock(ctx context.Context, tx *gorm.DB, bookID uint, quantity uint) error {
	book := new(models.Book)
	if err := tx.WithContext(ctx).First(book, "id = ?", bookID).Error; err != nil {
		return err
	}

	if book.Stock < quantity {
		return errors.New("book out of stock")
	}

	book.Stock = book.Stock - quantity

	if err := tx.WithContext(ctx).Where("id = ?", bookID).Updates(book).Error; err != nil {
		return err
	}

	return nil
}
