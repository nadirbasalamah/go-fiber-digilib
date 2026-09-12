package settings

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestGetAll_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	g := getall{repository: gormDB}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "settings" WHERE "settings"."deleted_at" IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "key", "value", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "test1", "test val", nil, nil, nil).
			AddRow(int64(2), "test2", "test val", nil, nil, nil))

	result, err := g.GetAll(context.Background())
	require.NoError(t, err)
	require.Equal(t, int(2), len(result))

	require.Equal(t, uint(1), result[0].ID)
	require.Equal(t, "test1", result[0].Key)
	require.Equal(t, uint(2), result[1].ID)
	require.Equal(t, "test2", result[1].Key)
}

func TestGetAll_ErrorOnFind(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	g := getall{repository: gormDB}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "settings" WHERE "settings"."deleted_at" IS NULL`)).
		WillReturnError(errors.New("database error"))

	_, err := g.GetAll(context.Background())
	require.Error(t, err)
	require.Contains(t, err.Error(), "database error")
}
