package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-digilib/api/middlewares"
	"go-digilib/carts"
	cartmocks "go-digilib/carts/mocks"
	"go-digilib/db/models"
	"go-digilib/pkg/dtos"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupCartsHandler(t *testing.T) (Carts, *cartmocks.MockService) {
	mockSvc := cartmocks.NewMockService(t)
	handler := NewCarts(mockSvc)
	return handler, mockSvc
}

var testUserClaims = &middlewares.JWTCustomClaims{
	ID:   1,
	Role: models.Enduser,
}

func TestCarts_GetByUser_Success(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupCartsHandler(t)

	result := []carts.Cart{
		{ID: 1, BookID: 1, UserID: 1, Quantity: 2},
		{ID: 2, BookID: 2, UserID: 1, Quantity: 1},
	}
	mockSvc.On("GetByUser", mock.Anything, uint(1)).Return(result, nil)

	req := httptest.NewRequest(http.MethodGet, "/carts/user", nil)
	rec := httptest.NewRecorder()
	e.GET("/carts/user", func(c *echo.Context) error {
		c.Set("userData", testUserClaims)
		return handler.GetByUser(c)
	})
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp dtos.Response[[]carts.Cart]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "success", resp.Status)
	require.Equal(t, "all carts", resp.Message)
	require.Len(t, resp.Data, 2)
}

func TestCarts_GetByUser_Error(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupCartsHandler(t)

	mockSvc.On("GetByUser", mock.Anything, uint(1)).Return([]carts.Cart{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/carts/user", nil)
	rec := httptest.NewRecorder()
	e.GET("/carts/user", func(c *echo.Context) error {
		c.Set("userData", testUserClaims)
		return handler.GetByUser(c)
	})
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "fetch carts failed", resp.Message)
}

func TestCarts_Create_Success(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupCartsHandler(t)

	cartReq := &carts.CartRequest{BookID: 1, Quantity: 2}
	result := carts.Cart{ID: 1, BookID: 1, UserID: 1, Quantity: 2}

	mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(req *carts.CartRequest) bool {
		return req.UserID == 1
	})).Return(result, nil)

	req := httptest.NewRequest(http.MethodPost, "/carts", strings.NewReader(`{"book_id":1,"quantity":2}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.POST("/carts", func(c *echo.Context) error {
		c.Set("userData", testUserClaims)
		c.Set("validatedBody", cartReq)
		return handler.Create(c)
	})
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	var resp dtos.Response[carts.Cart]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "success", resp.Status)
	require.Equal(t, "book added to the cart", resp.Message)
	require.Equal(t, uint(1), resp.Data.ID)
}

func TestCarts_Create_Error(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupCartsHandler(t)

	cartReq := &carts.CartRequest{BookID: 1, Quantity: 2}
	mockSvc.On("Create", mock.Anything, mock.Anything).Return(carts.Cart{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodPost, "/carts", strings.NewReader(`{"book_id":1,"quantity":2}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.POST("/carts", func(c *echo.Context) error {
		c.Set("userData", testUserClaims)
		c.Set("validatedBody", cartReq)
		return handler.Create(c)
	})
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "add book to the cart failed", resp.Message)
}

func TestCarts_Update_Success(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupCartsHandler(t)

	cartReq := &carts.CartRequest{BookID: 1, Quantity: 3}
	result := carts.Cart{ID: 1, BookID: 1, UserID: 1, Quantity: 3}

	mockSvc.On("Update", mock.Anything, mock.MatchedBy(func(req *carts.CartRequest) bool {
		return req.UserID == 1
	}), uint(1)).Return(result, nil)

	req := httptest.NewRequest(http.MethodPatch, "/carts/1", strings.NewReader(`{"book_id":1,"quantity":3}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.PATCH("/carts/:id", func(c *echo.Context) error {
		c.Set("userData", testUserClaims)
		c.Set("validatedBody", cartReq)
		return handler.Update(c)
	})
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp dtos.Response[carts.Cart]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "success", resp.Status)
	require.Equal(t, "cart updated", resp.Message)
	require.Equal(t, uint(3), resp.Data.Quantity)
}

func TestCarts_Update_InvalidID(t *testing.T) {
	e := echo.New()
	handler, _ := setupCartsHandler(t)

	req := httptest.NewRequest(http.MethodPatch, "/carts/abc", strings.NewReader(`{"book_id":1,"quantity":3}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.PATCH("/carts/:id", func(c *echo.Context) error {
		c.Set("userData", testUserClaims)
		return handler.Update(c)
	})
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "invalid id", resp.Message)
}

func TestCarts_Update_Error(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupCartsHandler(t)

	cartReq := &carts.CartRequest{BookID: 1, Quantity: 3}
	mockSvc.On("Update", mock.Anything, mock.Anything, uint(1)).Return(carts.Cart{}, errors.New("not found"))

	req := httptest.NewRequest(http.MethodPatch, "/carts/1", strings.NewReader(`{"book_id":1,"quantity":3}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.PATCH("/carts/:id", func(c *echo.Context) error {
		c.Set("userData", testUserClaims)
		c.Set("validatedBody", cartReq)
		return handler.Update(c)
	})
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "update cart failed", resp.Message)
}

func TestCarts_Delete_Success(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupCartsHandler(t)

	mockSvc.On("Delete", mock.Anything, uint(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/carts/1", nil)
	rec := httptest.NewRecorder()
	e.DELETE("/carts/:id", func(c *echo.Context) error {
		c.Set("userData", testUserClaims)
		return handler.Delete(c)
	})
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "success", resp.Status)
	require.Equal(t, "book removed from the cart", resp.Message)
}

func TestCarts_Delete_InvalidID(t *testing.T) {
	e := echo.New()
	handler, _ := setupCartsHandler(t)

	req := httptest.NewRequest(http.MethodDelete, "/carts/abc", nil)
	rec := httptest.NewRecorder()
	e.DELETE("/carts/:id", func(c *echo.Context) error {
		c.Set("userData", testUserClaims)
		return handler.Delete(c)
	})
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "invalid id", resp.Message)
}

func TestCarts_Delete_NotFound(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupCartsHandler(t)

	mockSvc.On("Delete", mock.Anything, uint(999)).Return(errors.New("not found"))

	req := httptest.NewRequest(http.MethodDelete, "/carts/999", nil)
	rec := httptest.NewRecorder()
	e.DELETE("/carts/:id", func(c *echo.Context) error {
		c.Set("userData", testUserClaims)
		return handler.Delete(c)
	})
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "cart not found", resp.Message)
}

func TestNewCarts(t *testing.T) {
	mockSvc := cartmocks.NewMockService(t)
	handler := NewCarts(mockSvc)
	require.NotNil(t, handler)
	require.NotNil(t, handler.carts)
}
