package books

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"go-digilib/pkg/utils"
)

func TestGetAll_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	g := getall{repository: gormDB}

	pagination := utils.Pagination{
		Limit:   10,
		Page:    1,
		Sort:    "id asc",
		Search:  "test",
		Keyword: "title",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "books"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(2)))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE title LIKE`)).
		WithArgs("test%", 10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "publisher", "year", "stock", "image_link", "category_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "Test Book 1", "Description 1", "Publisher 1", int64(2024), int64(10), "https://example.com/1.jpg", int64(1), nil, nil, nil).
			AddRow(int64(2), "Test Book 2", "Description 2", "Publisher 2", int64(2025), int64(5), "https://example.com/2.jpg", int64(2), nil, nil, nil))

	mock.ExpectQuery(`SELECT \* FROM "categories"`).
		WithArgs(int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description"}).
			AddRow(int64(1), "Fiction", "Fiction books").
			AddRow(int64(2), "Science", "Science books"))

	result, err := g.GetAll(context.Background(), pagination)
	require.NoError(t, err)
	require.Equal(t, int64(2), result.TotalRows)
	require.Equal(t, 1, result.TotalPages)

	rows, ok := result.Rows.([]Book)
	require.True(t, ok)
	require.Len(t, rows, 2)
	require.Equal(t, uint(1), rows[0].ID)
	require.Equal(t, uint(2), rows[1].ID)
	require.Equal(t, "Fiction", rows[0].Category.Name)
	require.Equal(t, "Science", rows[1].Category.Name)
}

func TestGetAll_ErrorOnFind(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	g := getall{repository: gormDB}

	pagination := utils.Pagination{
		Limit:   10,
		Page:    1,
		Sort:    "id asc",
		Search:  "test",
		Keyword: "title",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "books"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(0)))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE title LIKE`)).
		WithArgs("test%", 10).
		WillReturnError(errors.New("database error"))

	_, err := g.GetAll(context.Background(), pagination)
	require.Error(t, err)
	require.Contains(t, err.Error(), "database error")
}
