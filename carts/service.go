package carts

import (
	"context"

	"gorm.io/gorm"
)

type Service interface {
	GetByUser(ctx context.Context, userId uint) ([]Cart, error)
	GetByID(ctx context.Context, id uint) (Cart, error)
	Create(ctx context.Context, cartReq *CartRequest) (Cart, error)
	Update(ctx context.Context, cartReq *CartRequest, id uint) (Cart, error)
	Delete(ctx context.Context, id uint) error
}

type service struct {
	get
	getbyuser
	create
	update
	delete
}

var _ Service = (*service)(nil)

func New(repository *gorm.DB) Service {
	return service{
		get:       get{repository: repository},
		getbyuser: getbyuser{repository: repository},
		create:    create{repository: repository},
		update:    update{repository: repository, get: get{repository: repository}},
		delete:    delete{repository: repository, get: get{repository: repository}},
	}
}
