package settings

import (
	"context"

	"gorm.io/gorm"
)

type get struct {
	repository *gorm.DB
}

func (g get) GetByID(ctx context.Context, id uint) (Setting, error) {
	setting := new(Setting)

	if err := g.repository.WithContext(ctx).First(setting, "id = ?", id).Error; err != nil {
		return Setting{}, err
	}

	return *setting, nil
}
