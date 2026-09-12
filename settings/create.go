package settings

import (
	"context"
	"go-fiber-digilib/db/models"

	"gorm.io/gorm"
)

type create struct {
	repository *gorm.DB
}

func (c create) Create(ctx context.Context, settingReq *SettingRequest) (Setting, error) {
	setting := models.Setting{
		Key:   settingReq.Key,
		Value: settingReq.Value,
	}

	result := c.repository.WithContext(ctx).Create(&setting)
	record := new(Setting)

	if err := result.Error; err != nil {
		return Setting{}, err
	}

	if err := result.WithContext(ctx).Last(record).Error; err != nil {
		return Setting{}, err
	}

	return *record, nil
}
