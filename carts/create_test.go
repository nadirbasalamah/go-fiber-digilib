package carts

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = mockDB.Close() })

	gormDB, err := gorm.Open(
		postgres.New(postgres.Config{Conn: mockDB}),
		&gorm.Config{SkipDefaultTransaction: true},
	)
	require.NoError(t, err)

	return gormDB, mock
}

func testReq() *CartRequest {
	return &CartRequest{
		BookID:   1,
		Quantity: 2,
		UserID:   1,
	}
}

func TestCreate_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	c := create{repository: gormDB}
	req := testReq()

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "publisher", "year", "stock", "image_link", "category_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "Test Book", "Desc", "Pub", int64(2024), int64(10), "img.jpg", int64(1), nil, nil, nil))

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "carts"`)).
		WithArgs(int64(1), int64(1), int64(2), false, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(1)))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "carts"."id","carts"."book_id"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "book_id", "user_id", "quantity", "is_rented", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), int64(1), int64(1), int64(2), false, nil, nil, nil))

	mock.ExpectQuery(`SELECT \* FROM "books" WHERE`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "publisher", "year", "stock", "image_link", "category_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "Test Book", "Desc", "Pub", int64(2024), int64(10), "img.jpg", int64(1), nil, nil, nil))

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password", "address", "province_id", "city_id", "district_id", "profile_picture", "role", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "user1", "user@test.com", "", "", int64(0), int64(0), int64(0), "", "user", nil, nil, nil))

	mock.ExpectCommit()

	cart, err := c.Create(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, uint(1), cart.ID)
	require.Equal(t, uint(1), cart.BookID)
	require.Equal(t, uint(2), cart.Quantity)
}

func TestCreate_ErrorOnBookNotFound(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	c := create{repository: gormDB}
	req := testReq()

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnError(errors.New("book not found"))

	mock.ExpectRollback()

	_, err := c.Create(context.Background(), req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "book not found")
}

func TestCreate_ErrorOutOfStock(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	c := create{repository: gormDB}
	req := testReq()

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "publisher", "year", "stock", "image_link", "category_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "Test Book", "Desc", "Pub", int64(2024), int64(1), "img.jpg", int64(1), nil, nil, nil))

	mock.ExpectRollback()

	_, err := c.Create(context.Background(), req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "book out of stock")
}

func TestCreate_ErrorOnInsert(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	c := create{repository: gormDB}
	req := testReq()

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "publisher", "year", "stock", "image_link", "category_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "Test Book", "Desc", "Pub", int64(2024), int64(10), "img.jpg", int64(1), nil, nil, nil))

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "carts"`)).
		WithArgs(int64(1), int64(1), int64(2), false, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("insert error"))

	mock.ExpectRollback()

	_, err := c.Create(context.Background(), req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "insert error")
}

func TestCreate_ErrorOnLastQuery(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	c := create{repository: gormDB}
	req := testReq()

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "publisher", "year", "stock", "image_link", "category_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "Test Book", "Desc", "Pub", int64(2024), int64(10), "img.jpg", int64(1), nil, nil, nil))

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "carts"`)).
		WithArgs(int64(1), int64(1), int64(2), false, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(1)))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "carts"."id","carts"."book_id"`)).
		WillReturnError(errors.New("last query error"))

	mock.ExpectRollback()

	_, err := c.Create(context.Background(), req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "last query error")
}
