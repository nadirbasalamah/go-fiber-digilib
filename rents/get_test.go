package rents

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

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

func TestGetByID_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	g := get{repository: gormDB}

	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rents" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "quantity", "fee", "courier", "duration", "status", "return_time", "created_at", "updated_at", "returned_at", "deleted_at"}).
			AddRow(int64(1), int64(1), int64(3), float64(15000), "jne", int64(7), "pending", now, now, now, now, nil))

	rent, err := g.GetByID(context.Background(), 1)
	require.NoError(t, err)
	require.Equal(t, uint(1), rent.ID)
	require.Equal(t, uint(3), rent.Quantity)
	require.Equal(t, float64(15000), rent.Fee)
	require.Equal(t, "pending", rent.Status)
}

func TestGetByID_NotFound(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	g := get{repository: gormDB}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rents" WHERE id =`)).
		WithArgs(int64(999), 1).
		WillReturnError(errors.New("record not found"))

	_, err := g.GetByID(context.Background(), 999)
	require.Error(t, err)
	require.Contains(t, err.Error(), "record not found")
}
