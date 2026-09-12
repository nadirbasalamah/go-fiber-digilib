package books

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestUpdate_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	getOp := get{repository: gormDB}
	u := update{repository: gormDB, get: getOp}
	req := testReq()
	const id uint = 1

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "books" SET`)).
		WithArgs("Test Book", "Test Description", "Test Publisher", int64(2024), int64(10), "https://example.com/image.jpg", int64(1), sqlmock.AnyArg(), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "publisher", "year", "stock", "image_link", "category_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "Test Book", "Test Description", "Test Publisher", int64(2024), int64(10), "https://example.com/image.jpg", int64(1), nil, nil, nil))

	mock.ExpectQuery(`SELECT \* FROM "categories"`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description"}).
			AddRow(int64(1), "Fiction", "Fiction books"))

	book, err := u.Update(context.Background(), req, id)
	require.NoError(t, err)
	require.Equal(t, uint(1), book.ID)
	require.Equal(t, "Test Book", book.Title)
	require.Equal(t, "Fiction", book.Category.Name)
}

func TestUpdate_ErrorOnUpdate(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	getOp := get{repository: gormDB}
	u := update{repository: gormDB, get: getOp}
	req := testReq()
	const id uint = 1

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "books" SET`)).
		WithArgs("Test Book", "Test Description", "Test Publisher", int64(2024), int64(10), "https://example.com/image.jpg", int64(1), sqlmock.AnyArg(), int64(1)).
		WillReturnError(errors.New("update failed"))

	_, err := u.Update(context.Background(), req, id)
	require.Error(t, err)
	require.Contains(t, err.Error(), "update failed")
}

func TestUpdate_ErrorOnGetByID(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	getOp := get{repository: gormDB}
	u := update{repository: gormDB, get: getOp}
	req := testReq()
	const id uint = 1

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "books" SET`)).
		WithArgs("Test Book", "Test Description", "Test Publisher", int64(2024), int64(10), "https://example.com/image.jpg", int64(1), sqlmock.AnyArg(), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnError(errors.New("record not found"))

	_, err := u.Update(context.Background(), req, id)
	require.Error(t, err)
	require.Contains(t, err.Error(), "record not found")
}
