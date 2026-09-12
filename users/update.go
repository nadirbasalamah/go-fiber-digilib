package users

import (
	"context"
	"go-digilib/db/models"

	"gorm.io/gorm"
)

type update struct {
	repository *gorm.DB
	get
}

func (u update) Update(ctx context.Context, editReq *EditProfileRequest, id uint) (User, error) {
	user := models.User{
		Username:       editReq.Username,
		Email:          editReq.Email,
		Address:        editReq.Address,
		ProfilePicture: editReq.ProfilePicture,
	}

	result := u.repository.WithContext(ctx).Where("id = ?", id).Updates(&user)

	if err := result.Error; err != nil {
		return User{}, nil
	}

	record, err := u.get.GetProfile(ctx, id)

	if err != nil {
		return User{}, err
	}

	return record, nil
}
