package handlers

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"go-fiber-digilib/api/middlewares"
	"go-fiber-digilib/carts"
	cartsmocks "go-fiber-digilib/carts/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func testCartsHandler(t *testing.T) (Carts, *cartsmocks.MockService) {
	t.Helper()

	mockService := cartsmocks.NewMockService(t)

	return NewCarts(mockService), mockService
}

func userLocals(id int) map[string]any {
	return map[string]any{"userData": &middlewares.JWTCustomClaims{ID: id}}
}

func TestCartsGetByUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		handler, mockService := testCartsHandler(t)
		mockService.EXPECT().
			GetByUser(anyCtx(), uint(1)).
			Return([]carts.Cart{{ID: 1, UserID: 1, Quantity: 2}}, nil).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodGet, "/carts/user", handler.GetByUser, withLocals(userLocals(1)))

		resp := performRequest(t, app, http.MethodGet, "/carts/user", nil, "")
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "success", envelope.Status)
		require.Equal(t, "all carts", envelope.Message)
	})

	t.Run("service error", func(t *testing.T) {
		handler, mockService := testCartsHandler(t)
		mockService.EXPECT().
			GetByUser(anyCtx(), uint(1)).
			Return(nil, errors.New("db down")).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodGet, "/carts/user", handler.GetByUser, withLocals(userLocals(1)))

		resp := performRequest(t, app, http.MethodGet, "/carts/user", nil, "")
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "error occurred", envelope.Message)
	})
}

func TestCartsCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		handler, mockService := testCartsHandler(t)
		cartReq := &carts.CartRequest{BookID: 2, Quantity: 1}
		mockService.EXPECT().
			Create(anyCtx(), &carts.CartRequest{BookID: 2, Quantity: 1, UserID: 1}).
			Return(carts.Cart{ID: 1, BookID: 2, UserID: 1, Quantity: 1}, nil).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodPost, "/carts", handler.Create, withLocals(map[string]any{
			"userData":      &middlewares.JWTCustomClaims{ID: 1},
			"validatedBody": cartReq,
		}))

		resp := performRequest(t, app, http.MethodPost, "/carts", nil, "")
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "book added to the cart", envelope.Message)
	})

	t.Run("service error", func(t *testing.T) {
		handler, mockService := testCartsHandler(t)
		cartReq := &carts.CartRequest{BookID: 2, Quantity: 1}
		mockService.EXPECT().
			Create(anyCtx(), &carts.CartRequest{BookID: 2, Quantity: 1, UserID: 1}).
			Return(carts.Cart{}, errors.New("out of stock")).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodPost, "/carts", handler.Create, withLocals(map[string]any{
			"userData":      &middlewares.JWTCustomClaims{ID: 1},
			"validatedBody": cartReq,
		}))

		resp := performRequest(t, app, http.MethodPost, "/carts", nil, "")
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestCartsUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		handler, mockService := testCartsHandler(t)
		cartReq := &carts.CartRequest{BookID: 2, Quantity: 3}
		mockService.EXPECT().
			Update(anyCtx(), &carts.CartRequest{BookID: 2, Quantity: 3, UserID: 1}, uint(1)).
			Return(carts.Cart{ID: 1, Quantity: 3}, nil).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodPatch, "/carts/:id", handler.Update, withLocals(map[string]any{
			"userData":      &middlewares.JWTCustomClaims{ID: 1},
			"validatedBody": cartReq,
		}))

		resp := performRequest(t, app, http.MethodPatch, "/carts/1", nil, "")
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "cart updated", envelope.Message)
	})

	t.Run("invalid id", func(t *testing.T) {
		handler, _ := testCartsHandler(t)

		app := newTestApp(t)
		registerRoute(app, http.MethodPatch, "/carts/:id", handler.Update, withLocals(userLocals(1)))

		resp := performRequest(t, app, http.MethodPatch, "/carts/abc", nil, "")
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "invalid id", envelope.Message)
	})

	t.Run("service error", func(t *testing.T) {
		handler, mockService := testCartsHandler(t)
		cartReq := &carts.CartRequest{BookID: 2, Quantity: 3}
		mockService.EXPECT().
			Update(anyCtx(), &carts.CartRequest{BookID: 2, Quantity: 3, UserID: 1}, uint(1)).
			Return(carts.Cart{}, errors.New("not found")).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodPatch, "/carts/:id", handler.Update, withLocals(map[string]any{
			"userData":      &middlewares.JWTCustomClaims{ID: 1},
			"validatedBody": cartReq,
		}))

		resp := performRequest(t, app, http.MethodPatch, "/carts/1", nil, "")
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestCartsDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		handler, mockService := testCartsHandler(t)
		mockService.EXPECT().
			Delete(mock.MatchedBy(func(ctx context.Context) bool {
				return ctx.Value("userID") == 1
			}), uint(1)).
			Return(nil).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodDelete, "/carts/:id", handler.Delete, withLocals(userLocals(1)))

		resp := performRequest(t, app, http.MethodDelete, "/carts/1", nil, "")
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "book removed from the cart", envelope.Message)
	})

	t.Run("invalid id", func(t *testing.T) {
		handler, _ := testCartsHandler(t)

		app := newTestApp(t)
		registerRoute(app, http.MethodDelete, "/carts/:id", handler.Delete, withLocals(userLocals(1)))

		resp := performRequest(t, app, http.MethodDelete, "/carts/abc", nil, "")
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("service error", func(t *testing.T) {
		handler, mockService := testCartsHandler(t)
		mockService.EXPECT().
			Delete(anyCtx(), uint(1)).
			Return(errors.New("not found")).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodDelete, "/carts/:id", handler.Delete, withLocals(userLocals(1)))

		resp := performRequest(t, app, http.MethodDelete, "/carts/1", nil, "")
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
