package settings

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

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "settings" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "key", "value", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "setting key", "setting value", nil, nil, nil))

	setting, err := g.GetByID(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, uint(1), setting.ID)
	require.Equal(t, "setting key", setting.Key)
	require.Equal(t, "setting value", setting.Value)
}

func TestGetByID_NotFound(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	g := get{repository: gormDB}

	const id uint = 999

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "settings" WHERE id =`)).
		WithArgs(int64(999), 1).
		WillReturnError(errors.New("record not found"))

	_, err := g.GetByID(context.Background(), id)
	require.Error(t, err)
	require.Contains(t, err.Error(), "record not found")
}
