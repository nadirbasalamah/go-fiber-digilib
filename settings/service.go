package settings

import (
	"context"

	"gorm.io/gorm"
)

type Service interface {
	GetAll(ctx context.Context) ([]Setting, error)
	GetByID(ctx context.Context, id uint) (Setting, error)
	GetByKey(ctx context.Context, key string) (Setting, error)
	Create(ctx context.Context, settingReq *SettingRequest) (Setting, error)
	Update(ctx context.Context, settingReq *SettingRequest, id uint) (Setting, error)
	Delete(ctx context.Context, id uint) error
}

type service struct {
	getall
	get
	getbykey
	create
	update
	delete
}

var _ Service = (*service)(nil)

func New(repository *gorm.DB) Service {
	return service{
		getall:   getall{repository: repository},
		get:      get{repository: repository},
		getbykey: getbykey{repository: repository},
		create:   create{repository: repository},
		update:   update{repository: repository, get: get{repository: repository}},
		delete:   delete{repository: repository, get: get{repository: repository}},
	}
}
