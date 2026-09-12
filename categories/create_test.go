package categories

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

func testReq() *CategoryRequest {
	return &CategoryRequest{
		Name:        "Fiction",
		Description: "Fiction books",
	}
}

func TestCreate_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	c := create{repository: gormDB}
	req := testReq()

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "categories"`)).
		WithArgs("Fiction", "Fiction books", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(1)))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "categories"."id","categories"."name"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "Fiction", "Fiction books", nil, nil, nil))

	category, err := c.Create(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, uint(1), category.ID)
	require.Equal(t, "Fiction", category.Name)
	require.Equal(t, "Fiction books", category.Description)
}

func TestCreate_ErrorOnInsert(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	c := create{repository: gormDB}
	req := testReq()

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "categories"`)).
		WithArgs("Fiction", "Fiction books", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("insert error"))

	_, err := c.Create(context.Background(), req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "insert error")
}

func TestCreate_ErrorOnLastQuery(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	c := create{repository: gormDB}
	req := testReq()

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "categories"`)).
		WithArgs("Fiction", "Fiction books", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(1)))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "categories"."id","categories"."name"`)).
		WillReturnError(errors.New("last query error"))

	_, err := c.Create(context.Background(), req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "last query error")
}
