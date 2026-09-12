package settings

import (
	"context"
	"go-fiber-digilib/db/models"

	"gorm.io/gorm"
)

type update struct {
	repository *gorm.DB
	get
}

func (u update) Update(ctx context.Context, settingReq *SettingRequest, id uint) (Setting, error) {
	setting := models.Setting{
		Key:   settingReq.Key,
		Value: settingReq.Value,
	}

	result := u.repository.WithContext(ctx).Where("id = ?", id).Updates(&setting)

	if err := result.Error; err != nil {
		return Setting{}, err
	}

	record, err := u.get.GetByID(ctx, id)

	if err != nil {
		return Setting{}, err
	}

	return record, nil
}
