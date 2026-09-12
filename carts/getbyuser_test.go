package carts

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestGetByUser_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	g := getbyuser{repository: gormDB}

	const userID uint = 1

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "carts" WHERE user_id =`)).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "book_id", "user_id", "quantity", "is_rented", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), int64(1), int64(1), int64(2), false, nil, nil, nil).
			AddRow(int64(2), int64(2), int64(1), int64(1), false, nil, nil, nil))

	mock.ExpectQuery(`SELECT \* FROM "books" WHERE`).
		WithArgs(int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "publisher", "year", "stock", "image_link", "category_id", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "Book 1", "Desc 1", "Pub 1", int64(2024), int64(10), "img1.jpg", int64(1), nil, nil, nil).
			AddRow(int64(2), "Book 2", "Desc 2", "Pub 2", int64(2025), int64(5), "img2.jpg", int64(2), nil, nil, nil))

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password", "address", "province_id", "city_id", "district_id", "profile_picture", "role", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "user1", "user@test.com", "", "", int64(0), int64(0), int64(0), "", "user", nil, nil, nil))

	carts, err := g.GetByUser(context.Background(), userID)
	require.NoError(t, err)
	require.Len(t, carts, 2)
	require.Equal(t, uint(1), carts[0].ID)
	require.Equal(t, uint(2), carts[1].ID)
	require.Equal(t, uint(1), carts[0].Book.ID)
	require.Equal(t, uint(2), carts[1].Book.ID)
}

func TestGetByUser_Error(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	g := getbyuser{repository: gormDB}

	const userID uint = 1

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "carts" WHERE user_id =`)).
		WithArgs(int64(1)).
		WillReturnError(errors.New("database error"))

	_, err := g.GetByUser(context.Background(), userID)
	require.Error(t, err)
	require.Contains(t, err.Error(), "database error")
}
