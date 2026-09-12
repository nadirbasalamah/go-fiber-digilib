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

func TestGetByUser_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	g := getbyuser{repository: gormDB}

	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "user_rents"."id"`)).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "rent_id", "cart_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), int64(1), int64(1), now, now, nil).
			AddRow(int64(2), int64(2), int64(2), now, now, nil))

	mock.ExpectQuery(`SELECT \* FROM "carts" WHERE`).
		WithArgs(int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "book_id", "user_id", "quantity", "is_rented", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), int64(1), int64(1), int64(2), false, now, now, nil).
			AddRow(int64(2), int64(2), int64(1), int64(1), false, now, now, nil))

	mock.ExpectQuery(`SELECT \* FROM "books" WHERE`).
		WithArgs(int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "publisher", "year", "stock", "image_link", "category_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "Book1", "Desc", "Pub", int64(2024), int64(10), "img.jpg", int64(1), now, now, nil).
			AddRow(int64(2), "Book2", "Desc", "Pub", int64(2025), int64(5), "img2.jpg", int64(2), now, now, nil))

	mock.ExpectQuery(`SELECT \* FROM "categories" WHERE`).
		WithArgs(int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "Fiction", "Fiction books", now, now, nil).
			AddRow(int64(2), "Science", "Science books", now, now, nil))

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password", "address", "province_id", "city_id", "district_id", "profile_picture", "role", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "user1", "u@t.com", "", "", int64(0), int64(0), int64(0), "", "user", now, now, nil))

	mock.ExpectQuery(`SELECT \* FROM "rents" WHERE`).
		WithArgs(int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "quantity", "fee", "courier", "duration", "status", "return_time", "created_at", "updated_at", "returned_at", "deleted_at"}).
			AddRow(int64(1), int64(1), int64(3), float64(15000), "jne", int64(7), "pending", now, now, now, now, nil).
			AddRow(int64(2), int64(1), int64(1), float64(5000), "tiki", int64(7), "pending", now, now, now, now, nil))

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password", "address", "province_id", "city_id", "district_id", "profile_picture", "role", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "user1", "u@t.com", "", "", int64(0), int64(0), int64(0), "", "user", now, now, nil))

	result, err := g.GetByUser(context.Background(), 1)
	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Equal(t, uint(1), result[0].ID)
	require.Equal(t, uint(2), result[1].ID)
}

func TestGetByUser_Error(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	g := getbyuser{repository: gormDB}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "user_rents"."id"`)).
		WithArgs(int64(1)).
		WillReturnError(errors.New("database error"))

	_, err := g.GetByUser(context.Background(), 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "database error")
}
