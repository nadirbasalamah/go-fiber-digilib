package rents

import (
	"context"
	"go-fiber-digilib/pkg/utils"

	"gorm.io/gorm"
)

type Service interface {
	GetAll(ctx context.Context, pagination utils.Pagination) (utils.Pagination, error)
	GetByUser(ctx context.Context, userId uint) ([]UserRent, error)
	GetByID(ctx context.Context, id uint) (Rent, error)
	Create(ctx context.Context, rentReq *RentRequest) (Rent, error)
	Update(ctx context.Context, rentReq *RentUpdateRequest, id uint) (Rent, error)
	Delete(ctx context.Context, id uint) error
}

type service struct {
	getall
	get
	getbyuser
	create
	update
	delete
}

var _ Service = (*service)(nil)

func New(repository *gorm.DB) Service {
	return service{
		getall:    getall{repository: repository},
		get:       get{repository: repository},
		getbyuser: getbyuser{repository: repository},
		create:    create{repository: repository},
		update:    update{repository: repository, get: get{repository: repository}},
		delete:    delete{repository: repository, get: get{repository: repository}},
	}
}
