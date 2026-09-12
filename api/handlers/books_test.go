package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-digilib/books"
	bksmocks "go-digilib/books/mocks"
	"go-digilib/pkg/ai"
	"go-digilib/pkg/dtos"
	"go-digilib/pkg/utils"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupBooksHandler(t *testing.T) (Books, *bksmocks.MockService) {
	mockSvc := bksmocks.NewMockService(t)
	handler := NewBooks(mockSvc, nil, ai.Service{})
	return handler, mockSvc
}

func TestBooks_GetAll_Success(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupBooksHandler(t)

	pagination := utils.Pagination{Limit: 10, Page: 1}
	mockSvc.On("GetAll", mock.Anything, mock.Anything).Return(pagination, nil)

	req := httptest.NewRequest(http.MethodGet, "/books?page=1&limit=10&sort=id+asc&search=test", nil)
	rec := httptest.NewRecorder()
	e.GET("/books", func(c *echo.Context) error { return handler.GetAll(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp dtos.Response[utils.Pagination]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "success", resp.Status)
	require.Equal(t, "all books", resp.Message)
}

func TestBooks_GetAll_Error(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupBooksHandler(t)

	mockSvc.On("GetAll", mock.Anything, mock.Anything).Return(utils.Pagination{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	rec := httptest.NewRecorder()
	e.GET("/books", func(c *echo.Context) error { return handler.GetAll(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "fetch books failed", resp.Message)
}

func TestBooks_GetByID_Success(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupBooksHandler(t)

	book := books.Book{ID: 1, Title: "Test Book"}
	mockSvc.On("GetByID", mock.Anything, uint(1)).Return(book, nil)

	req := httptest.NewRequest(http.MethodGet, "/books/1", nil)
	rec := httptest.NewRecorder()
	e.GET("/books/:id", func(c *echo.Context) error { return handler.GetByID(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp dtos.Response[books.Book]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "success", resp.Status)
	require.Equal(t, "book found", resp.Message)
	require.Equal(t, uint(1), resp.Data.ID)
	require.Equal(t, "Test Book", resp.Data.Title)
}

func TestBooks_GetByID_InvalidID(t *testing.T) {
	e := echo.New()
	handler, _ := setupBooksHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/books/abc", nil)
	rec := httptest.NewRecorder()
	e.GET("/books/:id", func(c *echo.Context) error { return handler.GetByID(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "invalid id", resp.Message)
}

func TestBooks_GetByID_NotFound(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupBooksHandler(t)

	mockSvc.On("GetByID", mock.Anything, uint(999)).Return(books.Book{}, errors.New("not found"))

	req := httptest.NewRequest(http.MethodGet, "/books/999", nil)
	rec := httptest.NewRecorder()
	e.GET("/books/:id", func(c *echo.Context) error { return handler.GetByID(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "book not found", resp.Message)
}

func TestBooks_GetByCategory_Success(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupBooksHandler(t)

	pagination := utils.Pagination{Limit: 10, Page: 1}
	mockSvc.On("GetByCategory", mock.Anything, mock.Anything, uint(1)).Return(pagination, nil)

	req := httptest.NewRequest(http.MethodGet, "/books/category/1?page=1&limit=10", nil)
	rec := httptest.NewRecorder()
	e.GET("/books/category/:id", func(c *echo.Context) error { return handler.GetByCategory(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp dtos.Response[utils.Pagination]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "success", resp.Status)
	require.Equal(t, "all books by category", resp.Message)
}

func TestBooks_GetByCategory_InvalidID(t *testing.T) {
	e := echo.New()
	handler, _ := setupBooksHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/books/category/abc", nil)
	rec := httptest.NewRecorder()
	e.GET("/books/category/:id", func(c *echo.Context) error { return handler.GetByCategory(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "invalid id", resp.Message)
}

func TestBooks_GetByCategory_Error(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupBooksHandler(t)

	mockSvc.On("GetByCategory", mock.Anything, mock.Anything, uint(1)).Return(utils.Pagination{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/books/category/1", nil)
	rec := httptest.NewRecorder()
	e.GET("/books/category/:id", func(c *echo.Context) error { return handler.GetByCategory(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "fetch books failed", resp.Message)
}

func TestBooks_Delete_Success(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupBooksHandler(t)

	mockSvc.On("Delete", mock.Anything, uint(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/books/1", nil)
	rec := httptest.NewRecorder()
	e.DELETE("/books/:id", func(c *echo.Context) error { return handler.Delete(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "success", resp.Status)
	require.Equal(t, "book deleted", resp.Message)
}

func TestBooks_Delete_InvalidID(t *testing.T) {
	e := echo.New()
	handler, _ := setupBooksHandler(t)

	req := httptest.NewRequest(http.MethodDelete, "/books/abc", nil)
	rec := httptest.NewRecorder()
	e.DELETE("/books/:id", func(c *echo.Context) error { return handler.Delete(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "invalid id", resp.Message)
}

func TestBooks_Delete_NotFound(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupBooksHandler(t)

	mockSvc.On("Delete", mock.Anything, uint(999)).Return(errors.New("not found"))

	req := httptest.NewRequest(http.MethodDelete, "/books/999", nil)
	rec := httptest.NewRecorder()
	e.DELETE("/books/:id", func(c *echo.Context) error { return handler.Delete(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "book not found", resp.Message)
}

func TestNewBooks(t *testing.T) {
	mockSvc := bksmocks.NewMockService(t)
	handler := NewBooks(mockSvc, nil, ai.Service{})
	require.NotNil(t, handler)
	require.NotNil(t, handler.books)
}
