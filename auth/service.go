package auth

import (
	"context"

	"gorm.io/gorm"
)

type Service interface {
	Register(ctx context.Context, req *RegisterRequest) (User, error)
	Login(ctx context.Context, req *LoginRequest) (User, error)
}

type service struct {
	register
	login
}

var _ Service = (*service)(nil)

func New(repository *gorm.DB) Service {
	return service{
		register: register{repository: repository},
		login:    login{repository: repository},
	}
}
