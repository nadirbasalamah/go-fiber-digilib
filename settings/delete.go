package settings

import (
	"context"

	"gorm.io/gorm"
)

type delete struct {
	repository *gorm.DB
	get
}

func (d delete) Delete(ctx context.Context, id uint) error {
	setting, err := d.get.GetByID(ctx, id)

	if err != nil {
		return err
	}

	if err := d.repository.Delete(&setting).Error; err != nil {
		return err
	}

	return nil
}
