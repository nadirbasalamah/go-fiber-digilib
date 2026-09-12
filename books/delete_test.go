package books

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestDelete_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	getOp := get{repository: gormDB}
	d := delete{repository: gormDB, get: getOp}
	const id uint = 1

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "publisher", "year", "stock", "image_link", "category_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "Test Book", "Test Description", "Test Publisher", int64(2024), int64(10), "https://example.com/image.jpg", int64(1), nil, nil, nil))

	mock.ExpectQuery(`SELECT \* FROM "categories"`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description"}).
			AddRow(int64(1), "Fiction", "Fiction books"))

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "books" SET`)).
		WithArgs(sqlmock.AnyArg(), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := d.Delete(context.Background(), id)
	require.NoError(t, err)
}

func TestDelete_ErrorOnGetByID(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	getOp := get{repository: gormDB}
	d := delete{repository: gormDB, get: getOp}
	const id uint = 999

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE id =`)).
		WithArgs(int64(999), 1).
		WillReturnError(errors.New("record not found"))

	err := d.Delete(context.Background(), id)
	require.Error(t, err)
	require.Contains(t, err.Error(), "record not found")
}

func TestDelete_ErrorOnDelete(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	getOp := get{repository: gormDB}
	d := delete{repository: gormDB, get: getOp}
	const id uint = 1

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "publisher", "year", "stock", "image_link", "category_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "Test Book", "Test Description", "Test Publisher", int64(2024), int64(10), "https://example.com/image.jpg", int64(1), nil, nil, nil))

	mock.ExpectQuery(`SELECT \* FROM "categories"`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description"}).
			AddRow(int64(1), "Fiction", "Fiction books"))

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "books" SET`)).
		WithArgs(sqlmock.AnyArg(), int64(1)).
		WillReturnError(errors.New("delete failed"))

	err := d.Delete(context.Background(), id)
	require.Error(t, err)
	require.Contains(t, err.Error(), "delete failed")
}
