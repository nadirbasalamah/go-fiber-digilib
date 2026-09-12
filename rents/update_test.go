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

func TestUpdate_Success_BasicStatus(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	getOp := get{repository: gormDB}
	u := update{repository: gormDB, get: getOp}
	req := &RentUpdateRequest{Status: "rented"}

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "rents" SET`)).
		WithArgs("rented", sqlmock.AnyArg(), int64(1), "returned", "cancelled").
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rents" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "quantity", "fee", "courier", "duration", "status", "return_time", "created_at", "updated_at", "returned_at", "deleted_at"}).
			AddRow(int64(1), int64(1), int64(3), float64(15000), "jne", int64(7), "rented", now, now, now, now, nil))

	rent, err := u.Update(context.Background(), req, 1)
	require.NoError(t, err)
	require.Equal(t, "rented", rent.Status)
}

func TestUpdate_AlreadyReturned(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	getOp := get{repository: gormDB}
	u := update{repository: gormDB, get: getOp}
	req := &RentUpdateRequest{Status: "rented"}

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "rents" SET`)).
		WithArgs("rented", sqlmock.AnyArg(), int64(1), "returned", "cancelled").
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectRollback()

	_, err := u.Update(context.Background(), req, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "rent already returned or not found")
}

func TestUpdate_ErrorOnGetByID(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	getOp := get{repository: gormDB}
	u := update{repository: gormDB, get: getOp}
	req := &RentUpdateRequest{Status: "rented"}

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "rents" SET`)).
		WithArgs("rented", sqlmock.AnyArg(), int64(1), "returned", "cancelled").
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rents" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnError(errors.New("record not found"))

	_, err := u.Update(context.Background(), req, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "record not found")
}
