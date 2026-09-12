package settings

import (
	"context"

	"gorm.io/gorm"
)

type getall struct {
	repository *gorm.DB
}

func (g getall) GetAll(ctx context.Context) ([]Setting, error) {
	settings := []Setting{}

	if err := g.repository.WithContext(ctx).Find(&settings).Error; err != nil {
		return nil, err
	}

	return settings, nil
}
