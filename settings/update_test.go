package settings

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

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "settings" SET`)).
		WithArgs("setting key", "setting value", sqlmock.AnyArg(), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "settings" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "key", "value", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "setting key", "setting value", nil, nil, nil))

	setting, err := u.Update(context.Background(), req, id)
	require.NoError(t, err)
	require.Equal(t, uint(1), setting.ID)
	require.Equal(t, "setting key", setting.Key)
	require.Equal(t, "setting value", setting.Value)
}

func TestUpdate_ErrorOnUpdate(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	getOp := get{repository: gormDB}
	u := update{repository: gormDB, get: getOp}
	req := testReq()
	const id uint = 1

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "settings" SET`)).
		WithArgs("setting key", "setting value", sqlmock.AnyArg(), int64(1)).
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

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "settings" SET`)).
		WithArgs("setting key", "setting value", sqlmock.AnyArg(), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "settings" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnError(errors.New("record not found"))

	_, err := u.Update(context.Background(), req, id)
	require.Error(t, err)
	require.Contains(t, err.Error(), "record not found")
}
