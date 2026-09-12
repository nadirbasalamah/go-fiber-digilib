package settings

import (
	"context"
	"errors"
	"regexp"
	"testing"

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

func testReq() *SettingRequest {
	return &SettingRequest{
		Key:   "setting key",
		Value: "setting value",
	}
}

func TestCreate_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	c := create{repository: gormDB}
	req := testReq()

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "settings"`)).
		WithArgs("setting key", "setting value", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(1)))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "settings"."id","settings"."key"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "key", "value", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "setting key", "setting value", nil, nil, nil))

	setting, err := c.Create(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, uint(1), setting.ID)
	require.Equal(t, "setting key", setting.Key)
	require.Equal(t, "setting value", setting.Value)
}

func TestCreate_ErrorOnInsert(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	c := create{repository: gormDB}
	req := testReq()

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "settings"`)).
		WithArgs("setting key", "setting value", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("insert error"))

	_, err := c.Create(context.Background(), req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "insert error")
}

func TestCreate_ErrorOnLastQuery(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	c := create{repository: gormDB}
	req := testReq()

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "settings"`)).
		WithArgs("setting key", "setting value", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(1)))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "settings"."id","settings"."key"`)).
		WillReturnError(errors.New("last query error"))

	_, err := c.Create(context.Background(), req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "last query error")
}
