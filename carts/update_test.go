package carts

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

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "publisher", "year", "stock", "image_link", "category_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "Test Book", "Desc", "Pub", int64(2024), int64(10), "img.jpg", int64(1), nil, nil, nil))

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "carts" SET`)).
		WithArgs(int64(1), int64(2), sqlmock.AnyArg(), int64(1), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "carts" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "book_id", "user_id", "quantity", "is_rented", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), int64(1), int64(1), int64(2), false, nil, nil, nil))

	cart, err := u.Update(context.Background(), req, id)
	require.NoError(t, err)
	require.Equal(t, uint(1), cart.ID)
	require.Equal(t, uint(2), cart.Quantity)
}

func TestUpdate_ErrorOnBookNotFound(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	getOp := get{repository: gormDB}
	u := update{repository: gormDB, get: getOp}
	req := testReq()
	const id uint = 1

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnError(errors.New("book not found"))

	mock.ExpectRollback()

	_, err := u.Update(context.Background(), req, id)
	require.Error(t, err)
	require.Contains(t, err.Error(), "book not found")
}

func TestUpdate_ErrorOutOfStock(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	getOp := get{repository: gormDB}
	u := update{repository: gormDB, get: getOp}
	req := testReq()
	const id uint = 1

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "publisher", "year", "stock", "image_link", "category_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "Test Book", "Desc", "Pub", int64(2024), int64(1), "img.jpg", int64(1), nil, nil, nil))

	mock.ExpectRollback()

	_, err := u.Update(context.Background(), req, id)
	require.Error(t, err)
	require.Contains(t, err.Error(), "book out of stock")
}

func TestUpdate_ErrorOnUpdate(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	getOp := get{repository: gormDB}
	u := update{repository: gormDB, get: getOp}
	req := testReq()
	const id uint = 1

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "publisher", "year", "stock", "image_link", "category_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "Test Book", "Desc", "Pub", int64(2024), int64(10), "img.jpg", int64(1), nil, nil, nil))

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "carts" SET`)).
		WithArgs(int64(1), int64(2), sqlmock.AnyArg(), int64(1), int64(1)).
		WillReturnError(errors.New("update failed"))

	mock.ExpectRollback()

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

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "publisher", "year", "stock", "image_link", "category_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "Test Book", "Desc", "Pub", int64(2024), int64(10), "img.jpg", int64(1), nil, nil, nil))

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "carts" SET`)).
		WithArgs(int64(1), int64(2), sqlmock.AnyArg(), int64(1), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "carts" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnError(errors.New("record not found"))

	_, err := u.Update(context.Background(), req, id)
	require.Error(t, err)
	require.Contains(t, err.Error(), "record not found")
}
