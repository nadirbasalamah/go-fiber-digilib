package books

import (
	"context"
	"go-digilib/db/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type create struct {
	repository *gorm.DB
}

func (c create) Create(ctx context.Context, bookReq *BookRequest) (Book, error) {
	book := models.Book{
		Title:       bookReq.Title,
		Description: bookReq.Description,
		Publisher:   bookReq.Publisher,
		Year:        bookReq.Year,
		Stock:       bookReq.Stock,
		ImageLink:   bookReq.ImageLink,
		CategoryID:  bookReq.CategoryID,
	}

	result := c.repository.WithContext(ctx).Create(&book)
	record := new(Book)

	if err := result.Error; err != nil {
		return Book{}, err
	}

	if err := result.WithContext(ctx).Preload(clause.Associations).Last(record).Error; err != nil {
		return Book{}, err
	}

	return *record, nil
}
