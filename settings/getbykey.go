package settings

import (
	"context"

	"gorm.io/gorm"
)

type getbykey struct {
	repository *gorm.DB
}

func (g getbykey) GetByKey(ctx context.Context, key string) (Setting, error) {
	setting := new(Setting)

	if err := g.repository.WithContext(ctx).First(setting, "key = ?", key).Error; err != nil {
		return Setting{}, err
	}

	return *setting, nil
}
