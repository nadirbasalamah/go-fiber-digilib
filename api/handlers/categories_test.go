package handlers

import (
	"errors"
	"net/http"
	"testing"

	"go-fiber-digilib/categories"
	categoriesmocks "go-fiber-digilib/categories/mocks"
	"go-fiber-digilib/pkg/utils"

	"github.com/stretchr/testify/require"
)

func testCategoriesHandler(t *testing.T) (Categories, *categoriesmocks.MockService) {
	t.Helper()

	mockService := categoriesmocks.NewMockService(t)

	return NewCategories(mockService), mockService
}

func TestCategoriesGetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		handler, mockService := testCategoriesHandler(t)
		expected := utils.Pagination{TotalRows: 1, Rows: []categories.Category{{ID: 1, Name: "fiction"}}}

		mockService.EXPECT().
			GetAll(anyCtx(), utils.Pagination{Page: 1, Limit: 10, Sort: "id desc", Search: "fic", Keyword: "name"}).
			Return(expected, nil).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodGet, "/categories", handler.GetAll)

		resp := performRequest(t, app, http.MethodGet, "/categories?page=1&limit=10&sort=id+desc&search=fic", nil, "")
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "success", envelope.Status)
		require.Equal(t, "all categories", envelope.Message)
	})

	t.Run("service error", func(t *testing.T) {
		handler, mockService := testCategoriesHandler(t)
		mockService.EXPECT().
			GetAll(anyCtx(), utils.Pagination{Keyword: "name"}).
			Return(utils.Pagination{}, errors.New("db down")).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodGet, "/categories", handler.GetAll)

		resp := performRequest(t, app, http.MethodGet, "/categories", nil, "")
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "fetch categories failed", envelope.Message)
	})
}

func TestCategoriesGetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		handler, mockService := testCategoriesHandler(t)
		mockService.EXPECT().
			GetByID(anyCtx(), uint(1)).
			Return(categories.Category{ID: 1, Name: "fiction"}, nil).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodGet, "/categories/:id", handler.GetByID)

		resp := performRequest(t, app, http.MethodGet, "/categories/1", nil, "")
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "success", envelope.Status)
		require.Equal(t, "category found", envelope.Message)
	})

	t.Run("invalid id", func(t *testing.T) {
		handler, _ := testCategoriesHandler(t)

		app := newTestApp(t)
		registerRoute(app, http.MethodGet, "/categories/:id", handler.GetByID)

		resp := performRequest(t, app, http.MethodGet, "/categories/abc", nil, "")
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "invalid id", envelope.Message)
	})

	t.Run("service error", func(t *testing.T) {
		handler, mockService := testCategoriesHandler(t)
		mockService.EXPECT().
			GetByID(anyCtx(), uint(1)).
			Return(categories.Category{}, errors.New("not found")).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodGet, "/categories/:id", handler.GetByID)

		resp := performRequest(t, app, http.MethodGet, "/categories/1", nil, "")
		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "error occurred", envelope.Message)
	})
}

func TestCategoriesCreate(t *testing.T) {
	categoryReq := &categories.CategoryRequest{Name: "fiction", Description: "fiction books"}
	locals := map[string]any{"validatedBody": categoryReq}

	t.Run("success", func(t *testing.T) {
		handler, mockService := testCategoriesHandler(t)
		mockService.EXPECT().
			Create(anyCtx(), categoryReq).
			Return(categories.Category{ID: 1, Name: "fiction", Description: "fiction books"}, nil).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodPost, "/categories", handler.Create, withLocals(locals))

		resp := performRequest(t, app, http.MethodPost, "/categories", nil, "")
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "category created", envelope.Message)
	})

	t.Run("service error", func(t *testing.T) {
		handler, mockService := testCategoriesHandler(t)
		mockService.EXPECT().
			Create(anyCtx(), categoryReq).
			Return(categories.Category{}, errors.New("duplicate")).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodPost, "/categories", handler.Create, withLocals(locals))

		resp := performRequest(t, app, http.MethodPost, "/categories", nil, "")
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "error occurred", envelope.Message)
	})
}

func TestCategoriesUpdate(t *testing.T) {
	categoryReq := &categories.CategoryRequest{Name: "fiction", Description: "fiction books"}
	locals := map[string]any{"validatedBody": categoryReq}

	t.Run("success", func(t *testing.T) {
		handler, mockService := testCategoriesHandler(t)
		mockService.EXPECT().
			Update(anyCtx(), categoryReq, uint(1)).
			Return(categories.Category{ID: 1, Name: "fiction"}, nil).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodPatch, "/categories/:id", handler.Update, withLocals(locals))

		resp := performRequest(t, app, http.MethodPatch, "/categories/1", nil, "")
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "category updated", envelope.Message)
	})

	t.Run("invalid id", func(t *testing.T) {
		handler, _ := testCategoriesHandler(t)

		app := newTestApp(t)
		registerRoute(app, http.MethodPatch, "/categories/:id", handler.Update, withLocals(locals))

		resp := performRequest(t, app, http.MethodPatch, "/categories/abc", nil, "")
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "invalid id", envelope.Message)
	})

	t.Run("service error", func(t *testing.T) {
		handler, mockService := testCategoriesHandler(t)
		mockService.EXPECT().
			Update(anyCtx(), categoryReq, uint(1)).
			Return(categories.Category{}, errors.New("not found")).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodPatch, "/categories/:id", handler.Update, withLocals(locals))

		resp := performRequest(t, app, http.MethodPatch, "/categories/1", nil, "")
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestCategoriesDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		handler, mockService := testCategoriesHandler(t)
		mockService.EXPECT().
			Delete(anyCtx(), uint(1)).
			Return(nil).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodDelete, "/categories/:id", handler.Delete)

		resp := performRequest(t, app, http.MethodDelete, "/categories/1", nil, "")
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "category deleted", envelope.Message)
	})

	t.Run("invalid id", func(t *testing.T) {
		handler, _ := testCategoriesHandler(t)

		app := newTestApp(t)
		registerRoute(app, http.MethodDelete, "/categories/:id", handler.Delete)

		resp := performRequest(t, app, http.MethodDelete, "/categories/abc", nil, "")
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "invalid id", envelope.Message)
	})

	t.Run("service error", func(t *testing.T) {
		handler, mockService := testCategoriesHandler(t)
		mockService.EXPECT().
			Delete(anyCtx(), uint(1)).
			Return(errors.New("not found")).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodDelete, "/categories/:id", handler.Delete)

		resp := performRequest(t, app, http.MethodDelete, "/categories/1", nil, "")
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
