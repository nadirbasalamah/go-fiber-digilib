package carts

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

	ctx := context.WithValue(context.Background(), "userID", 1)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "carts" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "book_id", "user_id", "quantity", "is_rented", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), int64(1), int64(1), int64(2), false, nil, nil, nil))

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "carts" SET`)).
		WithArgs(sqlmock.AnyArg(), int64(1), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := d.Delete(ctx, id)
	require.NoError(t, err)
}

func TestDelete_ErrorOnGetByID(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	getOp := get{repository: gormDB}
	d := delete{repository: gormDB, get: getOp}
	const id uint = 999

	ctx := context.WithValue(context.Background(), "userID", 1)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "carts" WHERE id =`)).
		WithArgs(int64(999), 1).
		WillReturnError(errors.New("record not found"))

	err := d.Delete(ctx, id)
	require.Error(t, err)
	require.Contains(t, err.Error(), "record not found")
}

func TestDelete_ErrorOnDelete(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	getOp := get{repository: gormDB}
	d := delete{repository: gormDB, get: getOp}
	const id uint = 1

	ctx := context.WithValue(context.Background(), "userID", 1)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "carts" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "book_id", "user_id", "quantity", "is_rented", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), int64(1), int64(1), int64(2), false, nil, nil, nil))

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "carts" SET`)).
		WithArgs(sqlmock.AnyArg(), int64(1), int64(1)).
		WillReturnError(errors.New("delete failed"))

	err := d.Delete(ctx, id)
	require.Error(t, err)
	require.Contains(t, err.Error(), "delete cart failed")
}
