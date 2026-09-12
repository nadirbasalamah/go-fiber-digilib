package rents

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"go-fiber-digilib/pkg/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestGetAll_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	g := getall{repository: gormDB}

	pagination := utils.Pagination{
		Limit:   10,
		Page:    1,
		Sort:    "id asc",
		Search:  "",
		Keyword: "",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "user_rents"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(1)))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "user_rents"."id"`)).
		WithArgs("%", 10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "rent_id", "cart_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), int64(1), int64(1), nil, nil, nil))

	mock.ExpectQuery(`SELECT \* FROM "carts" WHERE`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "book_id", "user_id", "quantity", "is_rented", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), int64(1), int64(1), int64(2), false, nil, nil, nil))

	mock.ExpectQuery(`SELECT \* FROM "books" WHERE`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "publisher", "year", "stock", "image_link", "category_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "Book1", "Desc", "Pub", int64(2024), int64(10), "img.jpg", int64(1), nil, nil, nil))

	mock.ExpectQuery(`SELECT \* FROM "categories" WHERE`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "Fiction", "Fiction books", nil, nil, nil))

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password", "address", "province_id", "city_id", "district_id", "profile_picture", "role", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "user1", "u@t.com", "", "", int64(0), int64(0), int64(0), "", "user", nil, nil, nil))

	mock.ExpectQuery(`SELECT \* FROM "rents" WHERE`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "quantity", "fee", "courier", "duration", "status", "return_time", "created_at", "updated_at", "returned_at", "deleted_at"}).
			AddRow(int64(1), int64(1), int64(3), float64(15000), "jne", int64(7), "pending", time.Now(), time.Now(), time.Now(), time.Now(), nil))

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password", "address", "province_id", "city_id", "district_id", "profile_picture", "role", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "user1", "u@t.com", "", "", int64(0), int64(0), int64(0), "", "user", nil, nil, nil))

	result, err := g.GetAll(context.Background(), pagination)
	require.NoError(t, err)
	require.Equal(t, int64(1), result.TotalRows)

	rows, ok := result.Rows.([]UserRent)
	require.True(t, ok)
	require.Len(t, rows, 1)
	require.Equal(t, uint(1), rows[0].ID)
}

func TestGetAll_ErrorOnFind(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	g := getall{repository: gormDB}

	pagination := utils.Pagination{
		Limit:   10,
		Page:    1,
		Sort:    "id asc",
		Search:  "",
		Keyword: "",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "user_rents"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(0)))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "user_rents"."id"`)).
		WithArgs("%", 10).
		WillReturnError(errors.New("database error"))

	_, err := g.GetAll(context.Background(), pagination)
	require.Error(t, err)
	require.Contains(t, err.Error(), "database error")
}
