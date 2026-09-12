package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-digilib/api/middlewares"
	"go-digilib/pkg/dtos"
	"go-digilib/users"
	usermocks "go-digilib/users/mocks"

	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupUsersHandler(t *testing.T) (Users, *usermocks.MockService) {
	mockSvc := usermocks.NewMockService(t)
	handler := NewUsers(mockSvc, nil)
	return handler, mockSvc
}

func TestUsers_GetProfile_Unauthorized(t *testing.T) {
	e := echo.New()
	handler, _ := setupUsersHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	rec := httptest.NewRecorder()
	e.GET("/profile", func(c *echo.Context) error { return handler.GetProfile(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
}

func TestUsers_GetProfile_NotFound(t *testing.T) {
	e := echo.New()
	handler, mockSvc := setupUsersHandler(t)

	jwtCfg := middlewares.JWTConfig{SecretKey: "test-secret", ExpiresDuration: 60}
	token, err := jwtCfg.GenerateToken(1, "user")
	require.NoError(t, err)

	mockSvc.On("GetProfile", mock.Anything, uint(1)).Return(users.User{}, errors.New("not found"))

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	e.GET("/profile",
		middlewares.VerifyToken(func(c *echo.Context) error {
			return handler.GetProfile(c)
		}),
		echojwt.WithConfig(jwtCfg.Init()),
	)
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "user not found", resp.Message)
}

func TestUsers_EditProfile_InvalidToken(t *testing.T) {
	e := echo.New()
	handler, _ := setupUsersHandler(t)

	body := `{"username":"test","email":"test@test.com","address":"addr"}`
	req := httptest.NewRequest(http.MethodPatch, "/profile/edit", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.PATCH("/profile/edit", func(c *echo.Context) error { return handler.EditProfile(c) })
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
}

func TestUsers_EditProfile_InvalidBody(t *testing.T) {
	e := echo.New()
	handler, _ := setupUsersHandler(t)

	jwtCfg := middlewares.JWTConfig{SecretKey: "test-secret", ExpiresDuration: 60}
	token, err := jwtCfg.GenerateToken(1, "user")
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPatch, "/profile/edit", strings.NewReader(`invalid json`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	e.PATCH("/profile/edit",
		middlewares.VerifyToken(func(c *echo.Context) error {
			return handler.EditProfile(c)
		}),
		echojwt.WithConfig(jwtCfg.Init()),
	)
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "invalid request", resp.Message)
}

func TestUsers_EditProfile_ValidationError(t *testing.T) {
	e := echo.New()
	e.Validator = &middlewares.CustomValidator{Validator: middlewares.InitValidator()}
	handler, _ := setupUsersHandler(t)

	jwtCfg := middlewares.JWTConfig{SecretKey: "test-secret", ExpiresDuration: 60}
	token, err := jwtCfg.GenerateToken(1, "user")
	require.NoError(t, err)

	body := `{"username":"","email":"bad","address":""}`
	req := httptest.NewRequest(http.MethodPatch, "/profile/edit", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	e.PATCH("/profile/edit",
		middlewares.VerifyToken(func(c *echo.Context) error {
			return handler.EditProfile(c)
		}),
		echojwt.WithConfig(jwtCfg.Init()),
	)
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "validation failed", resp.Message)
}

func TestUsers_EditProfile_FileNotFound(t *testing.T) {
	e := echo.New()
	e.Validator = &middlewares.CustomValidator{Validator: middlewares.InitValidator()}
	handler, _ := setupUsersHandler(t)

	jwtCfg := middlewares.JWTConfig{SecretKey: "test-secret", ExpiresDuration: 60}
	token, err := jwtCfg.GenerateToken(1, "user")
	require.NoError(t, err)

	body := `{"username":"testuser","email":"test@test.com","address":"123 Test St"}`
	req := httptest.NewRequest(http.MethodPatch, "/profile/edit", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	e.PATCH("/profile/edit",
		middlewares.VerifyToken(func(c *echo.Context) error {
			return handler.EditProfile(c)
		}),
		echojwt.WithConfig(jwtCfg.Init()),
	)
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "file not found", resp.Message)
}

func TestNewUsers(t *testing.T) {
	mockSvc := usermocks.NewMockService(t)
	handler := NewUsers(mockSvc, nil)
	require.NotNil(t, handler)
	require.NotNil(t, handler.users)
}
