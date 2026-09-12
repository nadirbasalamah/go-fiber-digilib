package rents

import (
	"context"
	"errors"
	"go-fiber-digilib/db/models"

	"gorm.io/gorm"
)

type delete struct {
	repository *gorm.DB
	get
}

func (d delete) Delete(ctx context.Context, id uint) error {
	rent, err := d.get.GetByID(ctx, id)

	if err != nil {
		return err
	}

	isInvalidStatus := rent.Status != string(models.Returned) && rent.Status != string(models.Cancelled)

	if isInvalidStatus {
		return errors.New("rent status must be returned or cancelled")
	}

	err = d.repository.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).
			Delete(&models.UserRent{}, "rent_id = ?", id).
			Error; err != nil {
			return err
		}

		if err := tx.WithContext(ctx).Delete(&rent).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}
