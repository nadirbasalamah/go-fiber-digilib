package rents

import (
	"context"
	"fmt"
	"go-fiber-digilib/db/models"
	"time"

	"gorm.io/gorm"
)

type update struct {
	repository *gorm.DB
	get
}

func (u update) Update(ctx context.Context, rentReq *RentUpdateRequest, id uint) (Rent, error) {
	status := models.RentStatus(rentReq.Status)

	rent := models.Rent{
		Status: status,
	}

	if status == models.Returned {
		rent.ReturnedAt = time.Now()
	}

	txErr := u.repository.Transaction(func(tx *gorm.DB) error {
		result := tx.WithContext(ctx).
			Where("id = ? AND status NOT IN (?,?)", id, models.Returned, models.Cancelled).
			Updates(rent)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("rent already returned or not found")
		}

		if rent.Status == models.Returned || rent.Status == models.Cancelled {
			return returnBook(ctx, tx, id, rent)
		}

		return nil
	})

	if txErr != nil {
		return Rent{}, txErr
	}

	rentRecord, err := u.GetByID(ctx, id)

	if err != nil {
		return Rent{}, err
	}

	return rentRecord, nil
}

func returnBook(ctx context.Context, tx *gorm.DB, id uint, rent models.Rent) error {
	// get book IDs
	userRents := []models.UserRent{}

	if err := tx.WithContext(ctx).
		Where("rent_id = ?", id).
		Find(&userRents).
		Error; err != nil {
		return err
	}

	// update book stock
	for _, urent := range userRents {
		if err := returnBookStock(ctx, tx, urent.CartID); err != nil {
			return err
		}
	}

	// update rent data
	if err := tx.WithContext(ctx).
		Where("id = ?", id).
		Updates(&rent).
		Error; err != nil {
		return err
	}

	return nil
}

func returnBookStock(ctx context.Context, tx *gorm.DB, cartID uint) error {
	cart := new(models.Cart)

	if err := tx.WithContext(ctx).First(cart, "id = ?", cartID).Error; err != nil {
		return err
	}

	book := new(models.Book)

	if err := tx.WithContext(ctx).First(book, "id = ?", cart.BookID).Error; err != nil {
		return err
	}

	book.Stock = book.Stock + cart.Quantity

	if err := tx.WithContext(ctx).Where("id = ?", book.ID).Updates(book).Error; err != nil {
		return err
	}

	return nil
}
