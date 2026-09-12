package books

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestGetByID_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	g := get{repository: gormDB}

	const id uint = 1

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "publisher", "year", "stock", "image_link", "category_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "Test Book", "Test Description", "Test Publisher", int64(2024), int64(10), "https://example.com/image.jpg", int64(1), nil, nil, nil))

	mock.ExpectQuery(`SELECT \* FROM "categories"`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description"}).
			AddRow(int64(1), "Fiction", "Fiction books"))

	book, err := g.GetByID(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, uint(1), book.ID)
	require.Equal(t, "Test Book", book.Title)
	require.Equal(t, "Test Description", book.Description)
	require.Equal(t, "Test Publisher", book.Publisher)
	require.Equal(t, uint(2024), book.Year)
	require.Equal(t, uint(10), book.Stock)
	require.Equal(t, "https://example.com/image.jpg", book.ImageLink)
	require.Equal(t, uint(1), book.CategoryID)
	require.Equal(t, "Fiction", book.Category.Name)
	require.Equal(t, "Fiction books", book.Category.Description)
}

func TestGetByID_NotFound(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	g := get{repository: gormDB}

	const id uint = 999

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE id =`)).
		WithArgs(int64(999), 1).
		WillReturnError(errors.New("record not found"))

	_, err := g.GetByID(context.Background(), id)
	require.Error(t, err)
	require.Contains(t, err.Error(), "record not found")
}
