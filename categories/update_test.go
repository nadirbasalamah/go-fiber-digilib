package categories

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestUpdate_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	getOp := get{repository: gormDB}
	u := update{repository: gormDB, get: getOp}
	req := testReq()
	const id uint = 1

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "categories" SET`)).
		WithArgs("Fiction", "Fiction books", sqlmock.AnyArg(), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "Fiction", "Fiction books", nil, nil, nil))

	category, err := u.Update(context.Background(), req, id)
	require.NoError(t, err)
	require.Equal(t, uint(1), category.ID)
	require.Equal(t, "Fiction", category.Name)
	require.Equal(t, "Fiction books", category.Description)
}

func TestUpdate_ErrorOnUpdate(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	getOp := get{repository: gormDB}
	u := update{repository: gormDB, get: getOp}
	req := testReq()
	const id uint = 1

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "categories" SET`)).
		WithArgs("Fiction", "Fiction books", sqlmock.AnyArg(), int64(1)).
		WillReturnError(errors.New("update failed"))

	_, err := u.Update(context.Background(), req, id)
	require.Error(t, err)
	require.Contains(t, err.Error(), "update failed")
}

func TestUpdate_ErrorOnGetByID(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	getOp := get{repository: gormDB}
	u := update{repository: gormDB, get: getOp}
	req := testReq()
	const id uint = 1

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "categories" SET`)).
		WithArgs("Fiction", "Fiction books", sqlmock.AnyArg(), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnError(errors.New("record not found"))

	_, err := u.Update(context.Background(), req, id)
	require.Error(t, err)
	require.Contains(t, err.Error(), "record not found")
}
