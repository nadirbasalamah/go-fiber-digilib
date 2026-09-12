package users

import (
	"context"

	"gorm.io/gorm"
)

type Service interface {
	GetProfile(ctx context.Context, userID uint) (User, error)
	Update(ctx context.Context, editReq *EditProfileRequest, id uint) (User, error)
}

type service struct {
	get
	update
}

var _ Service = (*service)(nil)

func New(repository *gorm.DB) Service {
	return service{
		get:    get{repository: repository},
		update: update{repository: repository, get: get{repository: repository}},
	}
}
