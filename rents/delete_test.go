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

func TestDelete_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	getOp := get{repository: gormDB}
	d := delete{repository: gormDB, get: getOp}

	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rents" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "quantity", "fee", "courier", "duration", "status", "return_time", "created_at", "updated_at", "returned_at", "deleted_at"}).
			AddRow(int64(1), int64(1), int64(3), float64(15000), "jne", int64(7), "returned", now, now, now, now, nil))

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "user_rents" SET`)).
		WithArgs(sqlmock.AnyArg(), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 2))

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "rents" SET`)).
		WithArgs(sqlmock.AnyArg(), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	err := d.Delete(context.Background(), 1)
	require.NoError(t, err)
}

func TestDelete_ErrorOnGetByID(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	getOp := get{repository: gormDB}
	d := delete{repository: gormDB, get: getOp}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rents" WHERE id =`)).
		WithArgs(int64(999), 1).
		WillReturnError(errors.New("record not found"))

	err := d.Delete(context.Background(), 999)
	require.Error(t, err)
	require.Contains(t, err.Error(), "record not found")
}

func TestDelete_InvalidStatus(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	getOp := get{repository: gormDB}
	d := delete{repository: gormDB, get: getOp}

	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rents" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "quantity", "fee", "courier", "duration", "status", "return_time", "created_at", "updated_at", "returned_at", "deleted_at"}).
			AddRow(int64(1), int64(1), int64(3), float64(15000), "jne", int64(7), "pending", now, now, now, now, nil))

	err := d.Delete(context.Background(), 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "rent status must be returned or cancelled")
}

func TestDelete_ErrorInTransaction(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	getOp := get{repository: gormDB}
	d := delete{repository: gormDB, get: getOp}

	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rents" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "quantity", "fee", "courier", "duration", "status", "return_time", "created_at", "updated_at", "returned_at", "deleted_at"}).
			AddRow(int64(1), int64(1), int64(3), float64(15000), "jne", int64(7), "cancelled", now, now, now, now, nil))

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "user_rents" SET`)).
		WithArgs(sqlmock.AnyArg(), int64(1)).
		WillReturnError(errors.New("delete user_rents failed"))

	mock.ExpectRollback()

	err := d.Delete(context.Background(), 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "delete user_rents failed")
}
