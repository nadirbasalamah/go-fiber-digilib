package carts

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

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "carts" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "book_id", "user_id", "quantity", "is_rented", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), int64(1), int64(1), int64(2), false, nil, nil, nil))

	cart, err := g.GetByID(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, uint(1), cart.ID)
	require.Equal(t, uint(1), cart.BookID)
	require.Equal(t, uint(2), cart.Quantity)
}

func TestGetByID_NotFound(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	g := get{repository: gormDB}

	const id uint = 999

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "carts" WHERE id =`)).
		WithArgs(int64(999), 1).
		WillReturnError(errors.New("record not found"))

	_, err := g.GetByID(context.Background(), id)
	require.Error(t, err)
	require.Contains(t, err.Error(), "record not found")
}
