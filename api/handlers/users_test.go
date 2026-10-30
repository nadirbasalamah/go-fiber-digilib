package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-fiber-digilib/api/middlewares"
	"go-fiber-digilib/db/models"
	"go-fiber-digilib/pkg/fileupload"
	usersPkg "go-fiber-digilib/users"
	usersmocks "go-fiber-digilib/users/mocks"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUsersGetProfile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := usersmocks.NewMockService(t)
		handler := NewUsers(mockService, &stubUploader{})

		mockService.EXPECT().
			GetProfile(anyCtx(), uint(1)).
			Return(usersPkg.User{ID: 1, Username: "testuser"}, nil).
			Once()

		app, token := newAuthedApp(t, http.MethodGet, "/profile", handler.GetProfile)

		resp := performAuthed(t, app, http.MethodGet, "/profile", nil, "", token)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "user data", envelope.Message)
	})

	t.Run("invalid token", func(t *testing.T) {
		mockService := usersmocks.NewMockService(t)
		handler := NewUsers(mockService, &stubUploader{})

		app, _ := newAuthedApp(t, http.MethodGet, "/profile", handler.GetProfile)

		resp := performAuthed(t, app, http.MethodGet, "/profile", nil, "", "invalid.token.value")
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := usersmocks.NewMockService(t)
		handler := NewUsers(mockService, &stubUploader{})

		mockService.EXPECT().
			GetProfile(anyCtx(), uint(1)).
			Return(usersPkg.User{}, errors.New("not found")).
			Once()

		app, token := newAuthedApp(t, http.MethodGet, "/profile", handler.GetProfile)

		resp := performAuthed(t, app, http.MethodGet, "/profile", nil, "", token)
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

const testJWTSecret = "test-secret"

func newUsersHandler(t *testing.T, uploader fileupload.Uploader) (Users, *usersmocks.MockService) {
	t.Helper()

	mockService := usersmocks.NewMockService(t)

	return NewUsers(mockService, uploader), mockService
}

func newAuthedApp(t *testing.T, method, path string, handler fiber.Handler) (app *fiber.App, token string) {
	t.Helper()

	jwtConfig := middlewares.JWTConfig{SecretKey: testJWTSecret, ExpiresDuration: 60}

	token, err := jwtConfig.GenerateToken(1, models.Enduser)
	require.NoError(t, err)

	app = newTestApp(t)
	app.Add([]string{method}, path, jwtConfig.Init(), handler)

	return app, token
}

func performAuthed(t *testing.T, app *fiber.App, method, path string, body []byte, contentType, token string) *http.Response {
	t.Helper()

	req := httptest.NewRequest(method, path, readerFrom(body))
	req.Header.Set(fiber.HeaderAuthorization, "Bearer "+token)

	if contentType != "" {
		req.Header.Set(fiber.HeaderContentType, contentType)
	}

	resp, err := app.Test(req)
	require.NoError(t, err)

	return resp
}

func TestUsersEditProfile(t *testing.T) {
	validFields := map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
		"address":  "123 Test St",
	}

	t.Run("invalid token", func(t *testing.T) {
		mockService := usersmocks.NewMockService(t)
		handler := NewUsers(mockService, &stubUploader{})

		app, _ := newAuthedApp(t, http.MethodPatch, "/profile/edit", handler.EditProfile)

		resp := performAuthed(t, app, http.MethodPatch, "/profile/edit", nil, "", "invalid.token.value")
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("invalid form body", func(t *testing.T) {
		mockService := usersmocks.NewMockService(t)
		handler := NewUsers(mockService, &stubUploader{})

		app, token := newAuthedApp(t, http.MethodPatch, "/profile/edit", handler.EditProfile)

		resp := performAuthed(t, app, http.MethodPatch, "/profile/edit", []byte("not-a-form"), fiber.MIMEApplicationForm, token)
		require.NotEqual(t, http.StatusOK, resp.StatusCode)

		resp.Body.Close()
	})

	t.Run("validation failed", func(t *testing.T) {
		mockService := usersmocks.NewMockService(t)
		handler := NewUsers(mockService, &stubUploader{})

		body, contentType := newMultipartBody(t, map[string]string{"username": "", "email": "bad"}, "", "")

		app, token := newAuthedApp(t, http.MethodPatch, "/profile/edit", handler.EditProfile)

		resp := performAuthed(t, app, http.MethodPatch, "/profile/edit", body, contentType, token)
		require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "validation failed", envelope.Message)
	})

	t.Run("no file keeps old picture", func(t *testing.T) {
		mockService := usersmocks.NewMockService(t)
		handler := NewUsers(mockService, &stubUploader{})

		mockService.EXPECT().
			GetProfile(anyCtx(), uint(1)).
			Return(usersPkg.User{ID: 1, ProfilePicture: "https://img/old.png", ProfilePictureID: "old-id"}, nil).
			Once()
		mockService.EXPECT().
			Update(anyCtx(), mock.MatchedBy(func(req *usersPkg.EditProfileRequest) bool {
				return req.ProfilePicture == "https://img/old.png" && req.ProfilePictureID == "old-id"
			}), uint(1)).
			Return(usersPkg.User{ID: 1, Username: "testuser"}, nil).
			Once()

		body, contentType := newMultipartBody(t, validFields, "", "")

		app, token := newAuthedApp(t, http.MethodPatch, "/profile/edit", handler.EditProfile)

		resp := performAuthed(t, app, http.MethodPatch, "/profile/edit", body, contentType, token)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "profile updated", envelope.Message)
	})

	t.Run("invalid file format", func(t *testing.T) {
		mockService := usersmocks.NewMockService(t)
		handler := NewUsers(mockService, &stubUploader{})

		mockService.EXPECT().
			GetProfile(anyCtx(), uint(1)).
			Return(usersPkg.User{ID: 1, ProfilePictureID: "old-id"}, nil).
			Once()

		body, contentType := newMultipartBody(t, validFields, "avatar.txt", "hello")

		app, token := newAuthedApp(t, http.MethodPatch, "/profile/edit", handler.EditProfile)

		resp := performAuthed(t, app, http.MethodPatch, "/profile/edit", body, contentType, token)
		require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "invalid file format", envelope.Message)
	})

	t.Run("user not found", func(t *testing.T) {
		mockService := usersmocks.NewMockService(t)
		handler := NewUsers(mockService, &stubUploader{})

		mockService.EXPECT().
			GetProfile(anyCtx(), uint(1)).
			Return(usersPkg.User{}, errors.New("user not found")).
			Once()

		body, contentType := newMultipartBody(t, validFields, "avatar.png", "hello")

		app, token := newAuthedApp(t, http.MethodPatch, "/profile/edit", handler.EditProfile)

		resp := performAuthed(t, app, http.MethodPatch, "/profile/edit", body, contentType, token)
		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "user not found", envelope.Message)
	})

	t.Run("upload error", func(t *testing.T) {
		mockService := usersmocks.NewMockService(t)
		handler := NewUsers(mockService, &stubUploader{err: errors.New("cloudinary down")})

		mockService.EXPECT().
			GetProfile(anyCtx(), uint(1)).
			Return(usersPkg.User{ID: 1, ProfilePictureID: "old-id"}, nil).
			Once()

		body, contentType := newMultipartBody(t, validFields, "avatar.png", "hello")

		app, token := newAuthedApp(t, http.MethodPatch, "/profile/edit", handler.EditProfile)

		resp := performAuthed(t, app, http.MethodPatch, "/profile/edit", body, contentType, token)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "upload failed", envelope.Message)
	})

	t.Run("delete old file error", func(t *testing.T) {
		mockService := usersmocks.NewMockService(t)
		uploader := &stubUploader{url: "https://img/avatar.png", publicID: "new-id", deleteErr: errors.New("destroy failed")}
		handler := NewUsers(mockService, uploader)

		mockService.EXPECT().
			GetProfile(anyCtx(), uint(1)).
			Return(usersPkg.User{ID: 1, ProfilePictureID: "old-id"}, nil).
			Once()

		body, contentType := newMultipartBody(t, validFields, "avatar.png", "hello")

		app, token := newAuthedApp(t, http.MethodPatch, "/profile/edit", handler.EditProfile)

		resp := performAuthed(t, app, http.MethodPatch, "/profile/edit", body, contentType, token)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		require.Equal(t, "old-id", uploader.deletedID)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "old file failed to delete", envelope.Message)
	})

	t.Run("success", func(t *testing.T) {
		uploader := &stubUploader{url: "https://img/avatar.png", publicID: "new-id"}
		handler, mockService := newUsersHandler(t, uploader)

		mockService.EXPECT().
			GetProfile(anyCtx(), uint(1)).
			Return(usersPkg.User{ID: 1, ProfilePictureID: "old-id"}, nil).
			Once()
		mockService.EXPECT().
			Update(anyCtx(), mock.MatchedBy(func(req *usersPkg.EditProfileRequest) bool {
				return req.Username == "testuser" &&
					req.ProfilePicture == "https://img/avatar.png" &&
					req.ProfilePictureID == "new-id"
			}), uint(1)).
			Return(usersPkg.User{ID: 1, Username: "testuser", ProfilePicture: "https://img/avatar.png"}, nil).
			Once()

		body, contentType := newMultipartBody(t, validFields, "avatar.png", "hello")

		app, token := newAuthedApp(t, http.MethodPatch, "/profile/edit", handler.EditProfile)

		resp := performAuthed(t, app, http.MethodPatch, "/profile/edit", body, contentType, token)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "old-id", uploader.deletedID)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "profile updated", envelope.Message)
	})

	t.Run("success without file", func(t *testing.T) {
		handler, mockService := newUsersHandler(t, &stubUploader{})

		mockService.EXPECT().
			GetProfile(anyCtx(), uint(1)).
			Return(usersPkg.User{ID: 1, ProfilePicture: "https://img/old.png", ProfilePictureID: "old-id"}, nil).
			Once()
		mockService.EXPECT().
			Update(anyCtx(), mock.MatchedBy(func(req *usersPkg.EditProfileRequest) bool {
				return req.ProfilePicture == "https://img/old.png" && req.ProfilePictureID == "old-id"
			}), uint(1)).
			Return(usersPkg.User{ID: 1, Username: "testuser"}, nil).
			Once()

		body, contentType := newMultipartBody(t, validFields, "", "")

		app, token := newAuthedApp(t, http.MethodPatch, "/profile/edit", handler.EditProfile)

		resp := performAuthed(t, app, http.MethodPatch, "/profile/edit", body, contentType, token)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "profile updated", envelope.Message)
	})

	t.Run("service error", func(t *testing.T) {
		handler, mockService := newUsersHandler(t, &stubUploader{url: "https://img/avatar.png", publicID: "new-id"})

		mockService.EXPECT().
			GetProfile(anyCtx(), uint(1)).
			Return(usersPkg.User{ID: 1, ProfilePictureID: "old-id"}, nil).
			Once()
		mockService.EXPECT().
			Update(anyCtx(), mock.MatchedBy(func(req *usersPkg.EditProfileRequest) bool { return true }), uint(1)).
			Return(usersPkg.User{}, errors.New("db down")).
			Once()

		body, contentType := newMultipartBody(t, validFields, "avatar.png", "hello")

		app, token := newAuthedApp(t, http.MethodPatch, "/profile/edit", handler.EditProfile)

		resp := performAuthed(t, app, http.MethodPatch, "/profile/edit", body, contentType, token)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
