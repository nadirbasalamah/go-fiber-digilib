package handlers

import (
	"errors"
	"net/http"
	"testing"

	"go-fiber-digilib/settings"
	settingsmocks "go-fiber-digilib/settings/mocks"

	"github.com/stretchr/testify/require"
)

func testSettingsHandler(t *testing.T) (Settings, *settingsmocks.MockService) {
	t.Helper()

	mockService := settingsmocks.NewMockService(t)

	return NewSettings(mockService), mockService
}

func TestSettingsGetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		handler, mockService := testSettingsHandler(t)
		mockService.EXPECT().
			GetAll(anyCtx()).
			Return([]settings.Setting{{ID: 1, Key: "DISTRICT_ID", Value: "100"}}, nil).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodGet, "/settings", handler.GetAll)

		resp := performRequest(t, app, http.MethodGet, "/settings", nil, "")
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "all settings", envelope.Message)
	})

	t.Run("service error", func(t *testing.T) {
		handler, mockService := testSettingsHandler(t)
		mockService.EXPECT().
			GetAll(anyCtx()).
			Return(nil, errors.New("db down")).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodGet, "/settings", handler.GetAll)

		resp := performRequest(t, app, http.MethodGet, "/settings", nil, "")
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestSettingsGetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		handler, mockService := testSettingsHandler(t)
		mockService.EXPECT().
			GetByID(anyCtx(), uint(1)).
			Return(settings.Setting{ID: 1, Key: "DISTRICT_ID", Value: "100"}, nil).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodGet, "/settings/:id", handler.GetByID)

		resp := performRequest(t, app, http.MethodGet, "/settings/1", nil, "")
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "setting found", envelope.Message)
	})

	t.Run("invalid id", func(t *testing.T) {
		handler, _ := testSettingsHandler(t)

		app := newTestApp(t)
		registerRoute(app, http.MethodGet, "/settings/:id", handler.GetByID)

		resp := performRequest(t, app, http.MethodGet, "/settings/abc", nil, "")
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "invalid id", envelope.Message)
	})

	t.Run("service error", func(t *testing.T) {
		handler, mockService := testSettingsHandler(t)
		mockService.EXPECT().
			GetByID(anyCtx(), uint(1)).
			Return(settings.Setting{}, errors.New("not found")).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodGet, "/settings/:id", handler.GetByID)

		resp := performRequest(t, app, http.MethodGet, "/settings/1", nil, "")
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestSettingsGetByKey(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		handler, mockService := testSettingsHandler(t)
		mockService.EXPECT().
			GetByKey(anyCtx(), "DISTRICT_ID").
			Return(settings.Setting{ID: 1, Key: "DISTRICT_ID", Value: "100"}, nil).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodGet, "/settings/key/:key", handler.GetByKey)

		resp := performRequest(t, app, http.MethodGet, "/settings/key/DISTRICT_ID", nil, "")
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "setting found", envelope.Message)
	})

	t.Run("service error", func(t *testing.T) {
		handler, mockService := testSettingsHandler(t)
		mockService.EXPECT().
			GetByKey(anyCtx(), "UNKNOWN").
			Return(settings.Setting{}, errors.New("not found")).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodGet, "/settings/key/:key", handler.GetByKey)

		resp := performRequest(t, app, http.MethodGet, "/settings/key/UNKNOWN", nil, "")
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestSettingsCreate(t *testing.T) {
	settingReq := &settings.SettingRequest{Key: "DISTRICT_ID", Value: "100"}
	locals := map[string]any{"validatedBody": settingReq}

	t.Run("success", func(t *testing.T) {
		handler, mockService := testSettingsHandler(t)
		mockService.EXPECT().
			Create(anyCtx(), settingReq).
			Return(settings.Setting{ID: 1, Key: "DISTRICT_ID", Value: "100"}, nil).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodPost, "/settings", handler.Create, withLocals(locals))

		resp := performRequest(t, app, http.MethodPost, "/settings", nil, "")
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "setting created", envelope.Message)
	})

	t.Run("service error", func(t *testing.T) {
		handler, mockService := testSettingsHandler(t)
		mockService.EXPECT().
			Create(anyCtx(), settingReq).
			Return(settings.Setting{}, errors.New("duplicate")).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodPost, "/settings", handler.Create, withLocals(locals))

		resp := performRequest(t, app, http.MethodPost, "/settings", nil, "")
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestSettingsUpdate(t *testing.T) {
	settingReq := &settings.SettingRequest{Key: "DISTRICT_ID", Value: "200"}
	locals := map[string]any{"validatedBody": settingReq}

	t.Run("success", func(t *testing.T) {
		handler, mockService := testSettingsHandler(t)
		mockService.EXPECT().
			Update(anyCtx(), settingReq, uint(1)).
			Return(settings.Setting{ID: 1, Key: "DISTRICT_ID", Value: "200"}, nil).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodPatch, "/settings/:id", handler.Update, withLocals(locals))

		resp := performRequest(t, app, http.MethodPatch, "/settings/1", nil, "")
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "setting updated", envelope.Message)
	})

	t.Run("invalid id", func(t *testing.T) {
		handler, _ := testSettingsHandler(t)

		app := newTestApp(t)
		registerRoute(app, http.MethodPatch, "/settings/:id", handler.Update, withLocals(locals))

		resp := performRequest(t, app, http.MethodPatch, "/settings/abc", nil, "")
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("service error", func(t *testing.T) {
		handler, mockService := testSettingsHandler(t)
		mockService.EXPECT().
			Update(anyCtx(), settingReq, uint(1)).
			Return(settings.Setting{}, errors.New("not found")).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodPatch, "/settings/:id", handler.Update, withLocals(locals))

		resp := performRequest(t, app, http.MethodPatch, "/settings/1", nil, "")
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestSettingsDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		handler, mockService := testSettingsHandler(t)
		mockService.EXPECT().
			Delete(anyCtx(), uint(1)).
			Return(nil).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodDelete, "/settings/:id", handler.Delete)

		resp := performRequest(t, app, http.MethodDelete, "/settings/1", nil, "")
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "setting deleted", envelope.Message)
	})

	t.Run("invalid id", func(t *testing.T) {
		handler, _ := testSettingsHandler(t)

		app := newTestApp(t)
		registerRoute(app, http.MethodDelete, "/settings/:id", handler.Delete)

		resp := performRequest(t, app, http.MethodDelete, "/settings/abc", nil, "")
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("service error", func(t *testing.T) {
		handler, mockService := testSettingsHandler(t)
		mockService.EXPECT().
			Delete(anyCtx(), uint(1)).
			Return(errors.New("not found")).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodDelete, "/settings/:id", handler.Delete)

		resp := performRequest(t, app, http.MethodDelete, "/settings/1", nil, "")
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
