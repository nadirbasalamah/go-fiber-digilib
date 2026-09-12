package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-digilib/api/middlewares"
	"go-digilib/db/models"
	"go-digilib/pkg/dtos"
	"go-digilib/pkg/rajaongkir"
	"go-digilib/pkg/utils"
	"go-digilib/rents"
	rentmocks "go-digilib/rents/mocks"
	"go-digilib/settings"
	setmocks "go-digilib/settings/mocks"
	"go-digilib/users"
	usermocks "go-digilib/users/mocks"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupRentsHandler(t *testing.T) (Rents, *rentmocks.MockService, *setmocks.MockService, *usermocks.MockService) {
	rentsMock := rentmocks.NewMockService(t)
	settingsMock := setmocks.NewMockService(t)
	usersMock := usermocks.NewMockService(t)
	handler := NewRents(rentsMock, settingsMock, usersMock, rajaongkir.Service{})
	return handler, rentsMock, settingsMock, usersMock
}

var testRentUserClaims = &middlewares.JWTCustomClaims{
	ID:   1,
	Role: models.Enduser,
}

func TestRents_GetAll_Success(t *testing.T) {
	e := echo.New()
	handler, rentsMock, _, _ := setupRentsHandler(t)

	pagination := utils.Pagination{Limit: 10, Page: 1}
	rentsMock.On("GetAll", mock.Anything, mock.Anything).Return(pagination, nil)

	req := httptest.NewRequest(http.MethodGet, "/rents?page=1&limit=10", nil)
	rec := httptest.NewRecorder()
	e.GET("/rents", func(c *echo.Context) error { return handler.GetAll(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp dtos.Response[utils.Pagination]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "success", resp.Status)
	require.Equal(t, "all book rents", resp.Message)
}

func TestRents_GetAll_Error(t *testing.T) {
	e := echo.New()
	handler, rentsMock, _, _ := setupRentsHandler(t)

	rentsMock.On("GetAll", mock.Anything, mock.Anything).Return(utils.Pagination{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/rents", nil)
	rec := httptest.NewRecorder()
	e.GET("/rents", func(c *echo.Context) error { return handler.GetAll(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "fetch book rents failed", resp.Message)
}

func TestRents_GetByUser_Success(t *testing.T) {
	e := echo.New()
	handler, rentsMock, _, _ := setupRentsHandler(t)

	result := []rents.UserRent{{ID: 1, RentID: 1, CartID: 1}}
	rentsMock.On("GetByUser", mock.Anything, uint(1)).Return(result, nil)

	req := httptest.NewRequest(http.MethodGet, "/rents/user", nil)
	rec := httptest.NewRecorder()
	e.GET("/rents/user", func(c *echo.Context) error {
		c.Set("userData", testRentUserClaims)
		return handler.GetByUser(c)
	})
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp dtos.Response[[]rents.UserRent]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "success", resp.Status)
	require.Equal(t, "all rents", resp.Message)
	require.Len(t, resp.Data, 1)
}

func TestRents_GetByUser_Error(t *testing.T) {
	e := echo.New()
	handler, rentsMock, _, _ := setupRentsHandler(t)

	rentsMock.On("GetByUser", mock.Anything, uint(1)).Return([]rents.UserRent{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/rents/user", nil)
	rec := httptest.NewRecorder()
	e.GET("/rents/user", func(c *echo.Context) error {
		c.Set("userData", testRentUserClaims)
		return handler.GetByUser(c)
	})
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "fetch rents failed", resp.Message)
}

func TestRents_Create_UserNotFound(t *testing.T) {
	e := echo.New()
	handler, _, _, usersMock := setupRentsHandler(t)

	usersMock.On("GetProfile", mock.Anything, uint(1)).Return(users.User{}, errors.New("not found"))

	body := `{"cart_items":[1,2],"duration":7,"courier":"jne"}`
	req := httptest.NewRequest(http.MethodPost, "/rents", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.POST("/rents", func(c *echo.Context) error {
		c.Set("userData", testRentUserClaims)
		c.Set("validatedBody", &rents.RentRequest{
			CartItems: []uint{1, 2},
			Duration:  7,
			Courier:   "jne",
		})
		return handler.Create(c)
	})
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "failed to retrieve user data", resp.Message)
}

func TestRents_Create_SettingNotFound(t *testing.T) {
	e := echo.New()
	handler, _, settingsMock, usersMock := setupRentsHandler(t)

	usersMock.On("GetProfile", mock.Anything, uint(1)).Return(users.User{ID: 1, DistrictID: 100}, nil)
	settingsMock.On("GetByKey", mock.Anything, "DISTRICT_ID").Return(settings.Setting{}, errors.New("not found"))

	body := `{"cart_items":[1,2],"duration":7,"courier":"jne"}`
	req := httptest.NewRequest(http.MethodPost, "/rents", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.POST("/rents", func(c *echo.Context) error {
		c.Set("userData", testRentUserClaims)
		c.Set("validatedBody", &rents.RentRequest{
			CartItems: []uint{1, 2},
			Duration:  7,
			Courier:   "jne",
		})
		return handler.Create(c)
	})
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "failed to retrieve origin", resp.Message)
}

func TestRents_Update_Success(t *testing.T) {
	e := echo.New()
	handler, rentsMock, _, _ := setupRentsHandler(t)

	result := rents.Rent{ID: 1, Status: "rented"}
	rentsMock.On("Update", mock.Anything, mock.Anything, uint(1)).Return(result, nil)

	body := `{"status":"rented"}`
	req := httptest.NewRequest(http.MethodPatch, "/rents/1", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.PATCH("/rents/:id", func(c *echo.Context) error {
		c.Set("validatedBody", &rents.RentUpdateRequest{Status: "rented"})
		return handler.Update(c)
	})
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp dtos.Response[rents.Rent]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "success", resp.Status)
	require.Equal(t, "rent updated", resp.Message)
	require.Equal(t, "rented", resp.Data.Status)
}

func TestRents_Update_InvalidID(t *testing.T) {
	e := echo.New()
	handler, _, _, _ := setupRentsHandler(t)

	body := `{"status":"rented"}`
	req := httptest.NewRequest(http.MethodPatch, "/rents/abc", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.PATCH("/rents/:id", func(c *echo.Context) error {
		return handler.Update(c)
	})
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "invalid id", resp.Message)
}

func TestRents_Update_Error(t *testing.T) {
	e := echo.New()
	handler, rentsMock, _, _ := setupRentsHandler(t)

	rentsMock.On("Update", mock.Anything, mock.Anything, uint(1)).Return(rents.Rent{}, errors.New("not found"))

	body := `{"status":"rented"}`
	req := httptest.NewRequest(http.MethodPatch, "/rents/1", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.PATCH("/rents/:id", func(c *echo.Context) error {
		c.Set("validatedBody", &rents.RentUpdateRequest{Status: "rented"})
		return handler.Update(c)
	})
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "update rent failed", resp.Message)
}

func TestRents_Delete_Success(t *testing.T) {
	e := echo.New()
	handler, rentsMock, _, _ := setupRentsHandler(t)

	rentsMock.On("Delete", mock.Anything, uint(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/rents/1", nil)
	rec := httptest.NewRecorder()
	e.DELETE("/rents/:id", func(c *echo.Context) error { return handler.Delete(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "success", resp.Status)
	require.Equal(t, "book rent removed", resp.Message)
}

func TestRents_Delete_InvalidID(t *testing.T) {
	e := echo.New()
	handler, _, _, _ := setupRentsHandler(t)

	req := httptest.NewRequest(http.MethodDelete, "/rents/abc", nil)
	rec := httptest.NewRecorder()
	e.DELETE("/rents/:id", func(c *echo.Context) error { return handler.Delete(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "invalid id", resp.Message)
}

func TestRents_Delete_NotFound(t *testing.T) {
	e := echo.New()
	handler, rentsMock, _, _ := setupRentsHandler(t)

	rentsMock.On("Delete", mock.Anything, uint(999)).Return(errors.New("not found"))

	req := httptest.NewRequest(http.MethodDelete, "/rents/999", nil)
	rec := httptest.NewRecorder()
	e.DELETE("/rents/:id", func(c *echo.Context) error { return handler.Delete(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "rent not found", resp.Message)
}

func TestNewRents(t *testing.T) {
	handler := NewRents(nil, nil, nil, rajaongkir.Service{})
	require.NotNil(t, handler)
}
