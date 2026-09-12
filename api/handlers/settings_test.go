package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-digilib/pkg/dtos"
	"go-digilib/settings"
	setmocks "go-digilib/settings/mocks"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupSettingsHandler(t *testing.T) (Settings, *setmocks.MockService) {
	mockSvc := setmocks.NewMockService(t)
	handler := NewSettings(mockSvc)
	return handler, mockSvc
}

func TestSettings_GetAll_Success(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupSettingsHandler(t)

	result := []settings.Setting{
		{ID: 1, Key: "KEY_1", Value: "value1"},
		{ID: 2, Key: "KEY_2", Value: "value2"},
	}
	mockSvc.On("GetAll", mock.Anything).Return(result, nil)

	req := httptest.NewRequest(http.MethodGet, "/settings", nil)
	rec := httptest.NewRecorder()
	e.GET("/settings", func(c *echo.Context) error { return handler.GetAll(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp dtos.Response[[]settings.Setting]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "success", resp.Status)
	require.Equal(t, "all settings", resp.Message)
	require.Len(t, resp.Data, 2)
}

func TestSettings_GetAll_Error(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupSettingsHandler(t)

	mockSvc.On("GetAll", mock.Anything).Return([]settings.Setting{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/settings", nil)
	rec := httptest.NewRecorder()
	e.GET("/settings", func(c *echo.Context) error { return handler.GetAll(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "fetch settings failed", resp.Message)
}

func TestSettings_GetByID_Success(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupSettingsHandler(t)

	result := settings.Setting{ID: 1, Key: "KEY_1", Value: "value1"}
	mockSvc.On("GetByID", mock.Anything, uint(1)).Return(result, nil)

	req := httptest.NewRequest(http.MethodGet, "/settings/1", nil)
	rec := httptest.NewRecorder()
	e.GET("/settings/:id", func(c *echo.Context) error { return handler.GetByID(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp dtos.Response[settings.Setting]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "success", resp.Status)
	require.Equal(t, "setting found", resp.Message)
	require.Equal(t, uint(1), resp.Data.ID)
	require.Equal(t, "KEY_1", resp.Data.Key)
}

func TestSettings_GetByID_InvalidID(t *testing.T) {
	e := echo.New()
	handler, _ := setupSettingsHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/settings/abc", nil)
	rec := httptest.NewRecorder()
	e.GET("/settings/:id", func(c *echo.Context) error { return handler.GetByID(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "invalid id", resp.Message)
}

func TestSettings_GetByID_NotFound(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupSettingsHandler(t)

	mockSvc.On("GetByID", mock.Anything, uint(999)).Return(settings.Setting{}, errors.New("not found"))

	req := httptest.NewRequest(http.MethodGet, "/settings/999", nil)
	rec := httptest.NewRecorder()
	e.GET("/settings/:id", func(c *echo.Context) error { return handler.GetByID(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "setting not found", resp.Message)
}

func TestSettings_GetByKey_Success(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupSettingsHandler(t)

	result := settings.Setting{ID: 1, Key: "DISTRICT_ID", Value: "123"}
	mockSvc.On("GetByKey", mock.Anything, "DISTRICT_ID").Return(result, nil)

	req := httptest.NewRequest(http.MethodGet, "/settings/key/DISTRICT_ID", nil)
	rec := httptest.NewRecorder()
	e.GET("/settings/key/:key", func(c *echo.Context) error { return handler.GetByKey(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp dtos.Response[settings.Setting]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "success", resp.Status)
	require.Equal(t, "setting found", resp.Message)
	require.Equal(t, "DISTRICT_ID", resp.Data.Key)
	require.Equal(t, "123", resp.Data.Value)
}

func TestSettings_GetByKey_NotFound(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupSettingsHandler(t)

	mockSvc.On("GetByKey", mock.Anything, "UNKNOWN").Return(settings.Setting{}, errors.New("not found"))

	req := httptest.NewRequest(http.MethodGet, "/settings/key/UNKNOWN", nil)
	rec := httptest.NewRecorder()
	e.GET("/settings/key/:key", func(c *echo.Context) error { return handler.GetByKey(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "setting not found", resp.Message)
}

func TestSettings_Create_Success(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupSettingsHandler(t)

	settingReq := &settings.SettingRequest{Key: "NEW_KEY", Value: "new_value"}
	result := settings.Setting{ID: 1, Key: "NEW_KEY", Value: "new_value"}
	mockSvc.On("Create", mock.Anything, mock.Anything).Return(result, nil)

	body := `{"key":"NEW_KEY","value":"new_value"}`
	req := httptest.NewRequest(http.MethodPost, "/settings", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.POST("/settings", func(c *echo.Context) error {
		c.Set("validatedBody", settingReq)
		return handler.Create(c)
	})
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	var resp dtos.Response[settings.Setting]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "success", resp.Status)
	require.Equal(t, "setting created", resp.Message)
	require.Equal(t, uint(1), resp.Data.ID)
}

func TestSettings_Create_Error(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupSettingsHandler(t)

	settingReq := &settings.SettingRequest{Key: "NEW_KEY", Value: "new_value"}
	mockSvc.On("Create", mock.Anything, mock.Anything).Return(settings.Setting{}, errors.New("db error"))

	body := `{"key":"NEW_KEY","value":"new_value"}`
	req := httptest.NewRequest(http.MethodPost, "/settings", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.POST("/settings", func(c *echo.Context) error {
		c.Set("validatedBody", settingReq)
		return handler.Create(c)
	})
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "create setting failed", resp.Message)
}

func TestSettings_Update_Success(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupSettingsHandler(t)

	settingReq := &settings.SettingRequest{Key: "KEY_1", Value: "updated_value"}
	result := settings.Setting{ID: 1, Key: "KEY_1", Value: "updated_value"}
	mockSvc.On("Update", mock.Anything, mock.Anything, uint(1)).Return(result, nil)

	body := `{"key":"KEY_1","value":"updated_value"}`
	req := httptest.NewRequest(http.MethodPatch, "/settings/1", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.PATCH("/settings/:id", func(c *echo.Context) error {
		c.Set("validatedBody", settingReq)
		return handler.Update(c)
	})
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp dtos.Response[settings.Setting]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "success", resp.Status)
	require.Equal(t, "setting updated", resp.Message)
	require.Equal(t, "updated_value", resp.Data.Value)
}

func TestSettings_Update_InvalidID(t *testing.T) {
	e := echo.New()
	handler, _ := setupSettingsHandler(t)

	body := `{"key":"KEY_1","value":"updated_value"}`
	req := httptest.NewRequest(http.MethodPatch, "/settings/abc", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.PATCH("/settings/:id", func(c *echo.Context) error { return handler.Update(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "invalid id", resp.Message)
}

func TestSettings_Update_Error(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupSettingsHandler(t)

	settingReq := &settings.SettingRequest{Key: "KEY_1", Value: "updated_value"}
	mockSvc.On("Update", mock.Anything, mock.Anything, uint(1)).Return(settings.Setting{}, errors.New("not found"))

	body := `{"key":"KEY_1","value":"updated_value"}`
	req := httptest.NewRequest(http.MethodPatch, "/settings/1", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.PATCH("/settings/:id", func(c *echo.Context) error {
		c.Set("validatedBody", settingReq)
		return handler.Update(c)
	})
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "update setting failed", resp.Message)
}

func TestSettings_Delete_Success(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupSettingsHandler(t)

	mockSvc.On("Delete", mock.Anything, uint(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/settings/1", nil)
	rec := httptest.NewRecorder()
	e.DELETE("/settings/:id", func(c *echo.Context) error { return handler.Delete(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "success", resp.Status)
	require.Equal(t, "setting deleted", resp.Message)
}

func TestSettings_Delete_InvalidID(t *testing.T) {
	e := echo.New()
	handler, _ := setupSettingsHandler(t)

	req := httptest.NewRequest(http.MethodDelete, "/settings/abc", nil)
	rec := httptest.NewRecorder()
	e.DELETE("/settings/:id", func(c *echo.Context) error { return handler.Delete(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "invalid id", resp.Message)
}

func TestSettings_Delete_NotFound(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupSettingsHandler(t)

	mockSvc.On("Delete", mock.Anything, uint(999)).Return(errors.New("not found"))

	req := httptest.NewRequest(http.MethodDelete, "/settings/999", nil)
	rec := httptest.NewRecorder()
	e.DELETE("/settings/:id", func(c *echo.Context) error { return handler.Delete(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "setting not found", resp.Message)
}

func TestNewSettings(t *testing.T) {
	mockSvc := setmocks.NewMockService(t)
	handler := NewSettings(mockSvc)
	require.NotNil(t, handler)
	require.NotNil(t, handler.settings)
}
