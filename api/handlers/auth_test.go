package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-digilib/api/middlewares"
	"go-digilib/auth"
	authmocks "go-digilib/auth/mocks"
	"go-digilib/db/models"
	"go-digilib/pkg/dtos"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupEcho() *echo.Echo {
	e := echo.New()
	e.Validator = &middlewares.CustomValidator{Validator: middlewares.InitValidator()}
	return e
}

func setupJWTConfig() middlewares.JWTConfig {
	return middlewares.JWTConfig{
		SecretKey:       "test-secret-key",
		ExpiresDuration: 60,
	}
}

func setupAuthHandler(t *testing.T) (Auth, *authmocks.MockService) {
	mockSvc := authmocks.NewMockService(t)
	jwtCfg := setupJWTConfig()
	handler := NewAuth(mockSvc, jwtCfg)
	return handler, mockSvc
}

func TestRegister_Success(t *testing.T) {
	e := setupEcho()
	handler, mockSvc := setupAuthHandler(t)

	body := `{"username":"testuser","email":"test@example.com","password":"Password1!","address":"123 St","province_id":1,"city_id":2,"district_id":3}`

	mockSvc.On("Register", mock.Anything, mock.Anything).Return(auth.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "user",
	}, nil)

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Register(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, rec.Code)

	var resp dtos.Response[auth.User]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "success", resp.Status)
	require.Equal(t, "user registered", resp.Message)
	require.Equal(t, uint(1), resp.Data.ID)
	require.Equal(t, "testuser", resp.Data.Username)
}

func TestRegister_InvalidBody(t *testing.T) {
	e := setupEcho()
	handler, _ := setupAuthHandler(t)

	body := `invalid json`

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Register(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "invalid request", resp.Message)
}

func TestRegister_ValidationError(t *testing.T) {
	e := setupEcho()
	handler, _ := setupAuthHandler(t)

	body := `{"username":"","email":"bad","password":"weak","address":"","province_id":0,"city_id":0,"district_id":0}`

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Register(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "validation failed", resp.Message)
}

func TestRegister_ServiceError(t *testing.T) {
	e := setupEcho()
	handler, mockSvc := setupAuthHandler(t)

	body := `{"username":"testuser","email":"test@example.com","password":"Password1!","address":"123 St","province_id":1,"city_id":2,"district_id":3}`

	mockSvc.On("Register", mock.Anything, mock.Anything).Return(auth.User{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Register(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "user regsitration failed", resp.Message)
}

func TestLogin_Success(t *testing.T) {
	e := setupEcho()
	handler, mockSvc := setupAuthHandler(t)

	body := `{"email":"test@example.com","password":"Password1!"}`

	mockSvc.On("Login", mock.Anything, mock.Anything).Return(auth.User{
		ID:    1,
		Email: "test@example.com",
		Role:  "user",
	}, nil)

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Login(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp dtos.Response[string]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "success", resp.Status)
	require.Equal(t, "login success", resp.Message)
	require.NotEmpty(t, resp.Data)
}

func TestLogin_InvalidBody(t *testing.T) {
	e := setupEcho()
	handler, _ := setupAuthHandler(t)

	body := `invalid json`

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Login(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "invalid request", resp.Message)
}

func TestLogin_ValidationError(t *testing.T) {
	e := setupEcho()
	handler, _ := setupAuthHandler(t)

	body := `{"email":"bad","password":""}`

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Login(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "validation failed", resp.Message)
}

func TestLogin_ServiceError(t *testing.T) {
	e := setupEcho()
	handler, mockSvc := setupAuthHandler(t)

	body := `{"email":"test@example.com","password":"Password1!"}`

	mockSvc.On("Login", mock.Anything, mock.Anything).Return(auth.User{}, errors.New("wrong credentials"))

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Login(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp dtos.Response[any]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "failed", resp.Status)
	require.Equal(t, "user login failed", resp.Message)
}

func TestLogin_TokenGeneration(t *testing.T) {
	e := setupEcho()
	handler, mockSvc := setupAuthHandler(t)

	body := `{"email":"test@example.com","password":"Password1!"}`

	mockSvc.On("Login", mock.Anything, mock.Anything).Return(auth.User{
		ID:    1,
		Email: "test@example.com",
		Role:  "user",
	}, nil)

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Login(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp dtos.Response[string]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "success", resp.Status)

	require.NotEmpty(t, resp.Data)
	require.Contains(t, resp.Data, "eyJ")
}

func TestNewAuth(t *testing.T) {
	mockSvc := authmocks.NewMockService(t)
	jwtCfg := setupJWTConfig()

	authHandler := NewAuth(mockSvc, jwtCfg)
	require.NotNil(t, authHandler)
	require.NotNil(t, authHandler.auth)
	require.NotNil(t, authHandler.jwtConfig)
}

func TestLogin_JWTClaims(t *testing.T) {
	e := setupEcho()
	handler, mockSvc := setupAuthHandler(t)

	body := `{"email":"test@example.com","password":"Password1!"}`

	mockSvc.On("Login", mock.Anything, mock.Anything).Return(auth.User{
		ID:    42,
		Email: "test@example.com",
		Role:  string(models.Admin),
	}, nil)

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Login(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp dtos.Response[string]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.NotEmpty(t, resp.Data)
}
