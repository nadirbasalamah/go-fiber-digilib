package categories

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"go-digilib/pkg/utils"
)

func TestGetAll_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	g := getall{repository: gormDB}

	pagination := utils.Pagination{
		Limit:   10,
		Page:    1,
		Sort:    "id asc",
		Search:  "test",
		Keyword: "name",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "categories"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(2)))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" WHERE name LIKE`)).
		WithArgs("test%", 10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), "Fiction", "Fiction books", nil, nil, nil).
			AddRow(int64(2), "Science", "Science books", nil, nil, nil))

	result, err := g.GetAll(context.Background(), pagination)
	require.NoError(t, err)
	require.Equal(t, int64(2), result.TotalRows)
	require.Equal(t, 1, result.TotalPages)

	rows, ok := result.Rows.([]Category)
	require.True(t, ok)
	require.Len(t, rows, 2)
	require.Equal(t, uint(1), rows[0].ID)
	require.Equal(t, "Fiction", rows[0].Name)
	require.Equal(t, uint(2), rows[1].ID)
	require.Equal(t, "Science", rows[1].Name)
}

func TestGetAll_ErrorOnFind(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	g := getall{repository: gormDB}

	pagination := utils.Pagination{
		Limit:   10,
		Page:    1,
		Sort:    "id asc",
		Search:  "test",
		Keyword: "name",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "categories"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(0)))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" WHERE name LIKE`)).
		WithArgs("test%", 10).
		WillReturnError(errors.New("database error"))

	_, err := g.GetAll(context.Background(), pagination)
	require.Error(t, err)
	require.Contains(t, err.Error(), "database error")
}
