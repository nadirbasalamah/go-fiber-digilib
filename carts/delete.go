package carts

import (
	"context"
	"errors"
	"go-fiber-digilib/pkg/utils"

	"gorm.io/gorm"
)

type delete struct {
	repository *gorm.DB
	get
}

func (d delete) Delete(ctx context.Context, id uint) error {
	cart, err := d.get.GetByID(ctx, id)

	if err != nil {
		return err
	}

	userID := ctx.Value("userID").(int)

	res := d.repository.Scopes(utils.CurrentUser(uint(userID))).Delete(&cart)

	isFailed := res.Error != nil || res.RowsAffected == 0

	if isFailed {
		return errors.New("delete cart failed")
	}

	return nil
}
