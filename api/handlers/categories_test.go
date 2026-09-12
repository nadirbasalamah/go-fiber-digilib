package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-digilib/categories"
	catmocks "go-digilib/categories/mocks"
	"go-digilib/pkg/dtos"
	"go-digilib/pkg/utils"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupCategoriesHandler(t *testing.T) (Categories, *catmocks.MockService) {
	mockSvc := catmocks.NewMockService(t)
	handler := NewCategories(mockSvc)
	return handler, mockSvc
}

func TestCategories_GetAll_Success(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupCategoriesHandler(t)

	pagination := utils.Pagination{Limit: 10, Page: 1}
	mockSvc.On("GetAll", mock.Anything, mock.Anything).Return(pagination, nil)

	req := httptest.NewRequest(http.MethodGet, "/categories?page=1&limit=10&sort=id+asc&search=test", nil)
	rec := httptest.NewRecorder()
	e.GET("/categories", func(c *echo.Context) error { return handler.GetAll(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp dtos.Response[utils.Pagination]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "success", resp.Status)
	require.Equal(t, "all categories", resp.Message)
}

func TestCategories_GetAll_Error(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupCategoriesHandler(t)

	mockSvc.On("GetAll", mock.Anything, mock.Anything).Return(utils.Pagination{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	rec := httptest.NewRecorder()
	e.GET("/categories", func(c *echo.Context) error { return handler.GetAll(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "fetch categories failed", resp.Message)
}

func TestCategories_GetByID_Success(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupCategoriesHandler(t)

	category := categories.Category{ID: 1, Name: "Fiction", Description: "Fiction books"}
	mockSvc.On("GetByID", mock.Anything, uint(1)).Return(category, nil)

	req := httptest.NewRequest(http.MethodGet, "/categories/1", nil)
	rec := httptest.NewRecorder()
	e.GET("/categories/:id", func(c *echo.Context) error { return handler.GetByID(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp dtos.Response[categories.Category]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "success", resp.Status)
	require.Equal(t, "category found", resp.Message)
	require.Equal(t, uint(1), resp.Data.ID)
	require.Equal(t, "Fiction", resp.Data.Name)
}

func TestCategories_GetByID_InvalidID(t *testing.T) {
	e := echo.New()
	handler, _ := setupCategoriesHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/categories/abc", nil)
	rec := httptest.NewRecorder()
	e.GET("/categories/:id", func(c *echo.Context) error { return handler.GetByID(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "invalid id", resp.Message)
}

func TestCategories_GetByID_NotFound(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupCategoriesHandler(t)

	mockSvc.On("GetByID", mock.Anything, uint(999)).Return(categories.Category{}, errors.New("not found"))

	req := httptest.NewRequest(http.MethodGet, "/categories/999", nil)
	rec := httptest.NewRecorder()
	e.GET("/categories/:id", func(c *echo.Context) error { return handler.GetByID(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "category not found", resp.Message)
}

func TestCategories_Create_Success(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupCategoriesHandler(t)

	categoryReq := &categories.CategoryRequest{Name: "Fiction", Description: "Fiction books"}
	result := categories.Category{ID: 1, Name: "Fiction", Description: "Fiction books"}
	mockSvc.On("Create", mock.Anything, mock.Anything).Return(result, nil)

	body := `{"name":"Fiction","description":"Fiction books"}`
	req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.POST("/categories", func(c *echo.Context) error {
		c.Set("validatedBody", categoryReq)
		return handler.Create(c)
	})
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	var resp dtos.Response[categories.Category]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "success", resp.Status)
	require.Equal(t, "category created", resp.Message)
	require.Equal(t, uint(1), resp.Data.ID)
}

func TestCategories_Create_Error(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupCategoriesHandler(t)

	categoryReq := &categories.CategoryRequest{Name: "Fiction", Description: "Fiction books"}
	mockSvc.On("Create", mock.Anything, mock.Anything).Return(categories.Category{}, errors.New("db error"))

	body := `{"name":"Fiction","description":"Fiction books"}`
	req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.POST("/categories", func(c *echo.Context) error {
		c.Set("validatedBody", categoryReq)
		return handler.Create(c)
	})
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "create category failed", resp.Message)
}

func TestCategories_Update_Success(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupCategoriesHandler(t)

	categoryReq := &categories.CategoryRequest{Name: "Updated", Description: "Updated description"}
	result := categories.Category{ID: 1, Name: "Updated", Description: "Updated description"}
	mockSvc.On("Update", mock.Anything, mock.Anything, uint(1)).Return(result, nil)

	body := `{"name":"Updated","description":"Updated description"}`
	req := httptest.NewRequest(http.MethodPatch, "/categories/1", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.PATCH("/categories/:id", func(c *echo.Context) error {
		c.Set("validatedBody", categoryReq)
		return handler.Update(c)
	})
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp dtos.Response[categories.Category]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "success", resp.Status)
	require.Equal(t, "category updated", resp.Message)
	require.Equal(t, uint(1), resp.Data.ID)
}

func TestCategories_Update_InvalidID(t *testing.T) {
	e := echo.New()
	handler, _ := setupCategoriesHandler(t)

	body := `{"name":"Updated","description":"Updated description"}`
	req := httptest.NewRequest(http.MethodPatch, "/categories/abc", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.PATCH("/categories/:id", func(c *echo.Context) error { return handler.Update(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "invalid id", resp.Message)
}

func TestCategories_Update_Error(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupCategoriesHandler(t)

	categoryReq := &categories.CategoryRequest{Name: "Updated", Description: "Updated description"}
	mockSvc.On("Update", mock.Anything, mock.Anything, uint(1)).Return(categories.Category{}, errors.New("not found"))

	body := `{"name":"Updated","description":"Updated description"}`
	req := httptest.NewRequest(http.MethodPatch, "/categories/1", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.PATCH("/categories/:id", func(c *echo.Context) error {
		c.Set("validatedBody", categoryReq)
		return handler.Update(c)
	})
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "update category failed", resp.Message)
}

func TestCategories_Delete_Success(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupCategoriesHandler(t)

	mockSvc.On("Delete", mock.Anything, uint(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/categories/1", nil)
	rec := httptest.NewRecorder()
	e.DELETE("/categories/:id", func(c *echo.Context) error { return handler.Delete(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "success", resp.Status)
	require.Equal(t, "category deleted", resp.Message)
}

func TestCategories_Delete_InvalidID(t *testing.T) {
	e := echo.New()
	handler, _ := setupCategoriesHandler(t)

	req := httptest.NewRequest(http.MethodDelete, "/categories/abc", nil)
	rec := httptest.NewRecorder()
	e.DELETE("/categories/:id", func(c *echo.Context) error { return handler.Delete(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "invalid id", resp.Message)
}

func TestCategories_Delete_NotFound(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupCategoriesHandler(t)

	mockSvc.On("Delete", mock.Anything, uint(999)).Return(errors.New("not found"))

	req := httptest.NewRequest(http.MethodDelete, "/categories/999", nil)
	rec := httptest.NewRecorder()
	e.DELETE("/categories/:id", func(c *echo.Context) error { return handler.Delete(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "category not found", resp.Message)
}

func TestNewCategories(t *testing.T) {
	mockSvc := catmocks.NewMockService(t)
	handler := NewCategories(mockSvc)
	require.NotNil(t, handler)
	require.NotNil(t, handler.categories)
}
