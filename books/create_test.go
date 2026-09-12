package books

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

func testReq() *BookRequest {
	return &BookRequest{
		Title:       "Test Book",
		Description: "Test Description",
		Publisher:   "Test Publisher",
		Year:        2024,
		Stock:       10,
		CategoryID:  1,
		ImageLink:   "https://example.com/image.jpg",
	}
}

func TestCreate_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	c := create{repository: gormDB}
	req := testReq()

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "books"`)).
		WithArgs("Test Book", "Test Description", "Test Publisher", int64(2024), int64(10), "https://example.com/image.jpg", int64(1), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(1)))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "books"."id","books"."title"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "publisher", "year", "stock", "image_link", "category_id"}).
			AddRow(int64(1), "Test Book", "Test Description", "Test Publisher", int64(2024), int64(10), "https://example.com/image.jpg", int64(1)))

	mock.ExpectQuery(`SELECT \* FROM "categories"`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description"}).
			AddRow(int64(1), "Fiction", "Fiction books"))

	book, err := c.Create(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, uint(1), book.ID)
	require.Equal(t, "Test Book", book.Title)
	require.Equal(t, "Test Description", book.Description)
	require.Equal(t, "Test Publisher", book.Publisher)
	require.Equal(t, uint(2024), book.Year)
	require.Equal(t, uint(10), book.Stock)
	require.Equal(t, "https://example.com/image.jpg", book.ImageLink)
	require.Equal(t, uint(1), book.CategoryID)
	require.Equal(t, "Fiction", book.Category.Name)
	require.Equal(t, "Fiction books", book.Category.Description)
}

func TestCreate_ErrorOnInsert(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	c := create{repository: gormDB}
	req := testReq()

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "books"`)).
		WithArgs("Test Book", "Test Description", "Test Publisher", int64(2024), int64(10), "https://example.com/image.jpg", int64(1), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("insert error"))

	_, err := c.Create(context.Background(), req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "insert error")
}

func TestCreate_ErrorOnLastQuery(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	c := create{repository: gormDB}
	req := testReq()

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "books"`)).
		WithArgs("Test Book", "Test Description", "Test Publisher", int64(2024), int64(10), "https://example.com/image.jpg", int64(1), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(1)))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "books"."id","books"."title"`)).
		WillReturnError(errors.New("last query error"))

	_, err := c.Create(context.Background(), req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "last query error")
}

func TestCreate_ErrorOnPreload(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	c := create{repository: gormDB}
	req := testReq()

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "books"`)).
		WithArgs("Test Book", "Test Description", "Test Publisher", int64(2024), int64(10), "https://example.com/image.jpg", int64(1), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(1)))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "books"."id","books"."title"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "publisher", "year", "stock", "image_link", "category_id"}).
			AddRow(int64(1), "Test Book", "Test Description", "Test Publisher", int64(2024), int64(10), "https://example.com/image.jpg", int64(1)))

	mock.ExpectQuery(`SELECT \* FROM "categories"`).
		WithArgs(int64(1)).
		WillReturnError(errors.New("preload error"))

	_, err := c.Create(context.Background(), req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "preload error")
}

func TestCreate_EmptyRequest(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	c := create{repository: gormDB}
	req := &BookRequest{}

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "books"`)).
		WithArgs("", "", "", int64(0), int64(0), "", int64(0), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(1)))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "books"."id","books"."title"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "publisher", "year", "stock", "image_link", "category_id"}).
			AddRow(int64(1), "", "", "", int64(0), int64(0), "", int64(0)))

	book, err := c.Create(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, uint(1), book.ID)
	require.Empty(t, book.Title)
}
