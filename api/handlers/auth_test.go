package handlers

import (
	"errors"
	"net/http"
	"testing"

	"go-fiber-digilib/api/middlewares"
	"go-fiber-digilib/auth"
	authmocks "go-fiber-digilib/auth/mocks"
	"go-fiber-digilib/db/models"

	"github.com/stretchr/testify/require"
)

func testAuthHandler(t *testing.T) (Auth, *authmocks.MockService) {
	t.Helper()

	mockService := authmocks.NewMockService(t)

	return NewAuth(mockService, middlewares.JWTConfig{SecretKey: "secret", ExpiresDuration: 60}), mockService
}

func TestAuthRegister(t *testing.T) {
	validBody := map[string]any{
		"username":    "testuser",
		"email":       "test@example.com",
		"password":    "Password1!",
		"address":     "123 Test St",
		"province_id": 1,
		"city_id":     2,
		"district_id": 3,
	}

	t.Run("success", func(t *testing.T) {
		handler, mockService := testAuthHandler(t)
		mockService.EXPECT().
			Register(anyCtx(), &auth.RegisterRequest{
				Username:   "testuser",
				Email:      "test@example.com",
				Password:   "Password1!",
				Address:    "123 Test St",
				ProvinceID: 1,
				CityID:     2,
				DistrictID: 3,
			}).
			Return(auth.User{ID: 1, Username: "testuser", Email: "test@example.com"}, nil).
			Once()

		resp := performJSON(t, http.MethodPost, "/auth/register", jsonBody(t, validBody), handler.Register)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "success", envelope.Status)
		require.Equal(t, "user registered", envelope.Message)
	})

	t.Run("invalid json body", func(t *testing.T) {
		handler, _ := testAuthHandler(t)

		resp := performJSON(t, http.MethodPost, "/auth/register", []byte("{invalid"), handler.Register)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "failed", envelope.Status)
		require.Equal(t, "invalid request", envelope.Message)
	})

	t.Run("validation failed", func(t *testing.T) {
		handler, _ := testAuthHandler(t)

		body := map[string]any{"email": "not-an-email", "password": "weak"}

		resp := performJSON(t, http.MethodPost, "/auth/register", jsonBody(t, body), handler.Register)
		require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "validation failed", envelope.Message)
	})

	t.Run("service error", func(t *testing.T) {
		handler, mockService := testAuthHandler(t)
		mockService.EXPECT().
			Register(anyCtx(), &auth.RegisterRequest{
				Username:   "testuser",
				Email:      "test@example.com",
				Password:   "Password1!",
				Address:    "123 Test St",
				ProvinceID: 1,
				CityID:     2,
				DistrictID: 3,
			}).
			Return(auth.User{}, errors.New("duplicate email")).
			Once()

		resp := performJSON(t, http.MethodPost, "/auth/register", jsonBody(t, validBody), handler.Register)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "failed", envelope.Status)
		require.Equal(t, "error occurred", envelope.Message)
	})
}

func TestAuthLogin(t *testing.T) {
	validBody := map[string]any{
		"email":    "test@example.com",
		"password": "Password1!",
	}

	t.Run("success", func(t *testing.T) {
		handler, mockService := testAuthHandler(t)
		mockService.EXPECT().
			Login(anyCtx(), &auth.LoginRequest{Email: "test@example.com", Password: "Password1!"}).
			Return(auth.User{ID: 1, Email: "test@example.com", Role: string(models.Enduser)}, nil).
			Once()

		resp := performJSON(t, http.MethodPost, "/auth/login", jsonBody(t, validBody), handler.Login)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "success", envelope.Status)
		require.Equal(t, "login success", envelope.Message)

		var token string
		require.NoError(t, jsonUnmarshal(envelope.Data, &token))
		require.NotEmpty(t, token)
	})

	t.Run("invalid json body", func(t *testing.T) {
		handler, _ := testAuthHandler(t)

		resp := performJSON(t, http.MethodPost, "/auth/login", []byte("{invalid"), handler.Login)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "invalid request", envelope.Message)
	})

	t.Run("validation failed", func(t *testing.T) {
		handler, _ := testAuthHandler(t)

		resp := performJSON(t, http.MethodPost, "/auth/login", jsonBody(t, map[string]any{"email": "bad"}), handler.Login)
		require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "validation failed", envelope.Message)
	})

	t.Run("service error", func(t *testing.T) {
		handler, mockService := testAuthHandler(t)
		mockService.EXPECT().
			Login(anyCtx(), &auth.LoginRequest{Email: "test@example.com", Password: "Password1!"}).
			Return(auth.User{}, errors.New("invalid credentials")).
			Once()

		resp := performJSON(t, http.MethodPost, "/auth/login", jsonBody(t, validBody), handler.Login)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "failed", envelope.Status)
	})
}
