package handlers

import (
	"errors"
	"net/http"
	"testing"

	"go-fiber-digilib/pkg/utils"
	"go-fiber-digilib/rents"
	rentsmocks "go-fiber-digilib/rents/mocks"
	"go-fiber-digilib/settings"
	settingsmocks "go-fiber-digilib/settings/mocks"
	"go-fiber-digilib/users"
	usersmocks "go-fiber-digilib/users/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRentsGetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		rentsMock := rentsmocks.NewMockService(t)
		handler := NewRents(rentsMock, nil, nil, stubRajaOngkir{})
		expected := utils.Pagination{TotalRows: 1}

		rentsMock.EXPECT().
			GetAll(anyCtx(), utils.Pagination{Page: 1, Limit: 10, Keyword: "users.username"}).
			Return(expected, nil).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodGet, "/rents", handler.GetAll)

		resp := performRequest(t, app, http.MethodGet, "/rents?page=1&limit=10", nil, "")
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "all book rents", envelope.Message)
	})

	t.Run("service error", func(t *testing.T) {
		rentsMock := rentsmocks.NewMockService(t)
		handler := NewRents(rentsMock, nil, nil, stubRajaOngkir{})

		rentsMock.EXPECT().
			GetAll(anyCtx(), utils.Pagination{Keyword: "users.username"}).
			Return(utils.Pagination{}, errors.New("db down")).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodGet, "/rents", handler.GetAll)

		resp := performRequest(t, app, http.MethodGet, "/rents", nil, "")
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestRentsGetByUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		rentsMock := rentsmocks.NewMockService(t)
		handler := NewRents(rentsMock, nil, nil, stubRajaOngkir{})

		rentsMock.EXPECT().
			GetByUser(anyCtx(), uint(1)).
			Return([]rents.UserRent{{ID: 1, RentID: 2}}, nil).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodGet, "/rents/user", handler.GetByUser, withLocals(userLocals(1)))

		resp := performRequest(t, app, http.MethodGet, "/rents/user", nil, "")
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "all rents", envelope.Message)
	})

	t.Run("service error", func(t *testing.T) {
		rentsMock := rentsmocks.NewMockService(t)
		handler := NewRents(rentsMock, nil, nil, stubRajaOngkir{})

		rentsMock.EXPECT().
			GetByUser(anyCtx(), uint(1)).
			Return(nil, errors.New("db down")).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodGet, "/rents/user", handler.GetByUser, withLocals(userLocals(1)))

		resp := performRequest(t, app, http.MethodGet, "/rents/user", nil, "")
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestRentsCreate(t *testing.T) {
	rentReq := &rents.RentRequest{CartItems: []uint{1, 2}, Duration: 3, Courier: "jne"}
	locals := map[string]any{
		"userData":      userLocals(1)["userData"],
		"validatedBody": rentReq,
	}

	t.Run("user profile error", func(t *testing.T) {
		rentsMock := rentsmocks.NewMockService(t)
		usersMock := usersmocks.NewMockService(t)
		settingsMock := settingsmocks.NewMockService(t)
		handler := NewRents(rentsMock, settingsMock, usersMock, stubRajaOngkir{})

		usersMock.EXPECT().
			GetProfile(anyCtx(), uint(1)).
			Return(users.User{}, errors.New("not found")).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodPost, "/rents", handler.Create, withLocals(locals))

		resp := performRequest(t, app, http.MethodPost, "/rents", nil, "")
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "error occurred", envelope.Message)
	})

	t.Run("setting error", func(t *testing.T) {
		rentsMock := rentsmocks.NewMockService(t)
		usersMock := usersmocks.NewMockService(t)
		settingsMock := settingsmocks.NewMockService(t)
		handler := NewRents(rentsMock, settingsMock, usersMock, stubRajaOngkir{})

		usersMock.EXPECT().
			GetProfile(anyCtx(), uint(1)).
			Return(users.User{ID: 1, DistrictID: 10}, nil).
			Once()
		settingsMock.EXPECT().
			GetByKey(anyCtx(), "DISTRICT_ID").
			Return(settings.Setting{}, errors.New("not found")).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodPost, "/rents", handler.Create, withLocals(locals))

		resp := performRequest(t, app, http.MethodPost, "/rents", nil, "")
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "error occurred", envelope.Message)
	})

	t.Run("delivery fee error", func(t *testing.T) {
		rentsMock := rentsmocks.NewMockService(t)
		usersMock := usersmocks.NewMockService(t)
		settingsMock := settingsmocks.NewMockService(t)
		handler := NewRents(rentsMock, settingsMock, usersMock, stubRajaOngkir{err: errors.New("rajaongkir down")})

		usersMock.EXPECT().
			GetProfile(anyCtx(), uint(1)).
			Return(users.User{ID: 1, DistrictID: 10}, nil).
			Once()
		settingsMock.EXPECT().
			GetByKey(anyCtx(), "DISTRICT_ID").
			Return(settings.Setting{Key: "DISTRICT_ID", Value: "100"}, nil).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodPost, "/rents", handler.Create, withLocals(locals))

		resp := performRequest(t, app, http.MethodPost, "/rents", nil, "")
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "failed to calculate fee", envelope.Message)
	})

	t.Run("success", func(t *testing.T) {
		rentsMock := rentsmocks.NewMockService(t)
		usersMock := usersmocks.NewMockService(t)
		settingsMock := settingsmocks.NewMockService(t)
		handler := NewRents(rentsMock, settingsMock, usersMock, stubRajaOngkir{fee: 15000})

		usersMock.EXPECT().
			GetProfile(anyCtx(), uint(1)).
			Return(users.User{ID: 1, DistrictID: 10}, nil).
			Once()
		settingsMock.EXPECT().
			GetByKey(anyCtx(), "DISTRICT_ID").
			Return(settings.Setting{Key: "DISTRICT_ID", Value: "100"}, nil).
			Once()
		rentsMock.EXPECT().
			Create(anyCtx(), mock.MatchedBy(func(req *rents.RentRequest) bool {
				return req.Fee == 15000 && req.UserID == 1 && req.Courier == "jne" && !req.ReturnTime.IsZero()
			})).
			Return(rents.Rent{ID: 1, Fee: 15000, Courier: "jne"}, nil).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodPost, "/rents", handler.Create, withLocals(locals))

		resp := performRequest(t, app, http.MethodPost, "/rents", nil, "")
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "book added to the rent", envelope.Message)
	})

	t.Run("rent service error", func(t *testing.T) {
		rentsMock := rentsmocks.NewMockService(t)
		usersMock := usersmocks.NewMockService(t)
		settingsMock := settingsmocks.NewMockService(t)
		handler := NewRents(rentsMock, settingsMock, usersMock, stubRajaOngkir{fee: 15000})

		usersMock.EXPECT().
			GetProfile(anyCtx(), uint(1)).
			Return(users.User{ID: 1, DistrictID: 10}, nil).
			Once()
		settingsMock.EXPECT().
			GetByKey(anyCtx(), "DISTRICT_ID").
			Return(settings.Setting{Key: "DISTRICT_ID", Value: "100"}, nil).
			Once()
		rentsMock.EXPECT().
			Create(anyCtx(), mock.MatchedBy(func(req *rents.RentRequest) bool { return true })).
			Return(rents.Rent{}, errors.New("db down")).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodPost, "/rents", handler.Create, withLocals(locals))

		resp := performRequest(t, app, http.MethodPost, "/rents", nil, "")
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestRentsUpdate(t *testing.T) {
	updateReq := &rents.RentUpdateRequest{Status: "finished"}
	locals := map[string]any{"validatedBody": updateReq}

	t.Run("success", func(t *testing.T) {
		rentsMock := rentsmocks.NewMockService(t)
		handler := NewRents(rentsMock, nil, nil, stubRajaOngkir{})

		rentsMock.EXPECT().
			Update(anyCtx(), updateReq, uint(1)).
			Return(rents.Rent{ID: 1, Status: "finished"}, nil).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodPatch, "/rents/:id", handler.Update, withLocals(locals))

		resp := performRequest(t, app, http.MethodPatch, "/rents/1", nil, "")
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "rent updated", envelope.Message)
	})

	t.Run("invalid id", func(t *testing.T) {
		rentsMock := rentsmocks.NewMockService(t)
		handler := NewRents(rentsMock, nil, nil, stubRajaOngkir{})

		app := newTestApp(t)
		registerRoute(app, http.MethodPatch, "/rents/:id", handler.Update, withLocals(locals))

		resp := performRequest(t, app, http.MethodPatch, "/rents/abc", nil, "")
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "invalid id", envelope.Message)
	})

	t.Run("service error", func(t *testing.T) {
		rentsMock := rentsmocks.NewMockService(t)
		handler := NewRents(rentsMock, nil, nil, stubRajaOngkir{})

		rentsMock.EXPECT().
			Update(anyCtx(), updateReq, uint(1)).
			Return(rents.Rent{}, errors.New("not found")).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodPatch, "/rents/:id", handler.Update, withLocals(locals))

		resp := performRequest(t, app, http.MethodPatch, "/rents/1", nil, "")
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestRentsDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		rentsMock := rentsmocks.NewMockService(t)
		handler := NewRents(rentsMock, nil, nil, stubRajaOngkir{})

		rentsMock.EXPECT().
			Delete(anyCtx(), uint(1)).
			Return(nil).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodDelete, "/rents/:id", handler.Delete)

		resp := performRequest(t, app, http.MethodDelete, "/rents/1", nil, "")
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "book rent removed", envelope.Message)
	})

	t.Run("invalid id", func(t *testing.T) {
		rentsMock := rentsmocks.NewMockService(t)
		handler := NewRents(rentsMock, nil, nil, stubRajaOngkir{})

		app := newTestApp(t)
		registerRoute(app, http.MethodDelete, "/rents/:id", handler.Delete)

		resp := performRequest(t, app, http.MethodDelete, "/rents/abc", nil, "")
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("service error", func(t *testing.T) {
		rentsMock := rentsmocks.NewMockService(t)
		handler := NewRents(rentsMock, nil, nil, stubRajaOngkir{})

		rentsMock.EXPECT().
			Delete(anyCtx(), uint(1)).
			Return(errors.New("not found")).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodDelete, "/rents/:id", handler.Delete)

		resp := performRequest(t, app, http.MethodDelete, "/rents/1", nil, "")
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
