package rents

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func testRentReq() *RentRequest {
	now := time.Now()
	return &RentRequest{
		CartItems:  []uint{1, 2},
		Duration:   7,
		Courier:    "jne",
		Fee:        5000.0,
		UserID:     1,
		ReturnTime: now,
	}
}

func TestCreate_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	c := create{repository: gormDB}
	req := testRentReq()

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "carts" WHERE is_rented =`)).
		WithArgs(int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "book_id", "user_id", "quantity", "is_rented", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), int64(1), int64(1), int64(2), false, nil, nil, nil).
			AddRow(int64(2), int64(2), int64(1), int64(1), false, nil, nil, nil))

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "rents"`)).
		WithArgs(int64(1), int64(3), float64(15000), "jne", int64(7), "pending", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(1)))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rents"`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "quantity", "fee", "courier", "duration", "status", "return_time", "created_at", "updated_at", "returned_at", "deleted_at"}).
			AddRow(int64(1), int64(1), int64(3), float64(15000), "jne", int64(7), "pending", nil, nil, nil, nil, nil))

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password", "address", "province_id", "city_id", "district_id", "profile_picture", "role", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "user1", "u@t.com", "", "", int64(0), int64(0), int64(0), "", "user", nil, nil, nil))

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "user_rents"`)).
		WithArgs(int64(1), int64(1), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), int64(1), int64(2), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(1)))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "publisher", "year", "stock", "image_link", "category_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "Book1", "Desc", "Pub", int64(2024), int64(10), "img.jpg", int64(1), nil, nil, nil))

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "books" SET`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), int64(8), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), int64(1), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "carts" SET`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), true, sqlmock.AnyArg(), int64(1), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE id =`)).
		WithArgs(int64(2), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "publisher", "year", "stock", "image_link", "category_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(2), "Book2", "Desc", "Pub", int64(2025), int64(5), "img2.jpg", int64(2), nil, nil, nil))

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "books" SET`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), int64(4), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), int64(2), int64(2)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "carts" SET`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), true, sqlmock.AnyArg(), int64(2), int64(2)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	rent, err := c.Create(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, uint(1), rent.ID)
	require.Equal(t, uint(3), rent.Quantity)
	require.Equal(t, float64(15000), rent.Fee)
}

func TestCreate_CartsNotFound(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	c := create{repository: gormDB}
	req := testRentReq()

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "carts" WHERE is_rented =`)).
		WithArgs(int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "book_id", "user_id", "quantity", "is_rented", "created_at", "updated_at", "deleted_at"}))

	mock.ExpectRollback()

	_, err := c.Create(context.Background(), req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "carts not found")
}

func TestCreate_ErrorOnInsertRent(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	c := create{repository: gormDB}
	req := testRentReq()

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "carts" WHERE is_rented =`)).
		WithArgs(int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "book_id", "user_id", "quantity", "is_rented", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), int64(1), int64(1), int64(2), false, nil, nil, nil).
			AddRow(int64(2), int64(2), int64(1), int64(1), false, nil, nil, nil))

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "rents"`)).
		WithArgs(int64(1), int64(3), float64(15000), "jne", int64(7), "pending", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("insert failed"))

	mock.ExpectRollback()

	_, err := c.Create(context.Background(), req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "insert failed")
}
