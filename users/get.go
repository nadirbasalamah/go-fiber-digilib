package users

import (
	"context"

	"gorm.io/gorm"
)

type get struct {
	repository *gorm.DB
}

func (g get) GetProfile(ctx context.Context, userID uint) (User, error) {
	user := new(User)

	err := g.repository.WithContext(ctx).First(user, "id = ?", userID).Error

	if err != nil {
		return User{}, err
	}

	return *user, nil
}
