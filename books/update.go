package books

import (
	"context"
	"go-fiber-digilib/db/models"

	"gorm.io/gorm"
)

type update struct {
	repository *gorm.DB
	get
}

func (u update) Update(ctx context.Context, bookReq *BookRequest, id uint) (Book, error) {
	book := models.Book{
		Title:       bookReq.Title,
		Description: bookReq.Description,
		Publisher:   bookReq.Publisher,
		Year:        bookReq.Year,
		Stock:       bookReq.Stock,
		CategoryID:  bookReq.CategoryID,
	}

	if bookReq.ImageLink != "" {
		book.ImageLink = bookReq.ImageLink
		book.ImageID = bookReq.ImageID
	}

	result := u.repository.WithContext(ctx).Where("id = ?", id).Updates(&book)

	if err := result.Error; err != nil {
		return Book{}, err
	}

	record, err := u.get.GetByID(ctx, id)

	if err != nil {
		return Book{}, err
	}

	return record, nil
}
