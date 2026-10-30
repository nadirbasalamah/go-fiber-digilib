package handlers

import (
	"errors"
	"net/http"
	"testing"

	"go-fiber-digilib/books"
	booksmocks "go-fiber-digilib/books/mocks"
	"go-fiber-digilib/pkg/ai"
	"go-fiber-digilib/pkg/utils"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func testBooksHandler(t *testing.T) (Books, *booksmocks.MockService) {
	t.Helper()

	mockService := booksmocks.NewMockService(t)

	return NewBooks(mockService, stubUploader{}, stubRecommendation{}), mockService
}

func TestBooksGetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		handler, mockService := testBooksHandler(t)
		expected := utils.Pagination{TotalRows: 1, Rows: []books.Book{{ID: 1, Title: "Go"}}}

		mockService.EXPECT().
			GetAll(anyCtx(), utils.Pagination{Page: 2, Limit: 5, Sort: "title asc", Search: "go", Keyword: "title"}).
			Return(expected, nil).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodGet, "/books", handler.GetAll)

		resp := performRequest(t, app, http.MethodGet, "/books?page=2&limit=5&sort=title+asc&search=go", nil, "")
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "all books", envelope.Message)
	})

	t.Run("service error", func(t *testing.T) {
		handler, mockService := testBooksHandler(t)
		mockService.EXPECT().
			GetAll(anyCtx(), utils.Pagination{Keyword: "title"}).
			Return(utils.Pagination{}, errors.New("db down")).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodGet, "/books", handler.GetAll)

		resp := performRequest(t, app, http.MethodGet, "/books", nil, "")
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestBooksGetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		handler, mockService := testBooksHandler(t)
		mockService.EXPECT().
			GetByID(anyCtx(), uint(1)).
			Return(books.Book{ID: 1, Title: "Go"}, nil).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodGet, "/books/:id", handler.GetByID)

		resp := performRequest(t, app, http.MethodGet, "/books/1", nil, "")
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "book found", envelope.Message)
	})

	t.Run("invalid id", func(t *testing.T) {
		handler, _ := testBooksHandler(t)

		app := newTestApp(t)
		registerRoute(app, http.MethodGet, "/books/:id", handler.GetByID)

		resp := performRequest(t, app, http.MethodGet, "/books/abc", nil, "")
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "invalid id", envelope.Message)
	})

	t.Run("service error", func(t *testing.T) {
		handler, mockService := testBooksHandler(t)
		mockService.EXPECT().
			GetByID(anyCtx(), uint(1)).
			Return(books.Book{}, errors.New("not found")).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodGet, "/books/:id", handler.GetByID)

		resp := performRequest(t, app, http.MethodGet, "/books/1", nil, "")
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestBooksGetByCategory(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		handler, mockService := testBooksHandler(t)
		expected := utils.Pagination{TotalRows: 1}

		mockService.EXPECT().
			GetByCategory(anyCtx(), utils.Pagination{Page: 1, Limit: 10, Keyword: "title"}, uint(7)).
			Return(expected, nil).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodGet, "/books/category/:id", handler.GetByCategory)

		resp := performRequest(t, app, http.MethodGet, "/books/category/7?page=1&limit=10", nil, "")
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "all books by category", envelope.Message)
	})

	t.Run("invalid id", func(t *testing.T) {
		handler, _ := testBooksHandler(t)

		app := newTestApp(t)
		registerRoute(app, http.MethodGet, "/books/category/:id", handler.GetByCategory)

		resp := performRequest(t, app, http.MethodGet, "/books/category/abc", nil, "")
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("service error", func(t *testing.T) {
		handler, mockService := testBooksHandler(t)
		mockService.EXPECT().
			GetByCategory(anyCtx(), utils.Pagination{Keyword: "title"}, uint(7)).
			Return(utils.Pagination{}, errors.New("db down")).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodGet, "/books/category/:id", handler.GetByCategory)

		resp := performRequest(t, app, http.MethodGet, "/books/category/7", nil, "")
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestBooksCreate(t *testing.T) {
	bookFields := map[string]string{
		"title":       "Go",
		"description": "Golang book",
		"publisher":   "Publisher",
		"year":        "2024",
		"stock":       "3",
		"category_id": "1",
	}
	bookReq := &books.BookRequest{
		Title:       "Go",
		Description: "Golang book",
		Publisher:   "Publisher",
		Year:        2024,
		Stock:       3,
		CategoryID:  1,
	}
	locals := map[string]any{"validatedBody": bookReq}

	t.Run("file not found", func(t *testing.T) {
		handler, _ := testBooksHandler(t)

		body, contentType := newMultipartBody(t, bookFields, "", "")

		app := newTestApp(t)
		registerRoute(app, http.MethodPost, "/books", handler.Create, withLocals(locals))

		resp := performRequest(t, app, http.MethodPost, "/books", body, contentType)
		require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "file not found", envelope.Message)
	})

	t.Run("invalid file format", func(t *testing.T) {
		handler, _ := testBooksHandler(t)

		body, contentType := newMultipartBody(t, bookFields, "book.pdf", "content")

		app := newTestApp(t)
		registerRoute(app, http.MethodPost, "/books", handler.Create, withLocals(locals))

		resp := performRequest(t, app, http.MethodPost, "/books", body, contentType)
		require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "invalid file format", envelope.Message)
	})

	t.Run("upload error", func(t *testing.T) {
		mockService := booksmocks.NewMockService(t)
		handler := NewBooks(mockService, stubUploader{err: errors.New("cloudinary down")}, stubRecommendation{})

		body, contentType := newMultipartBody(t, bookFields, "book.png", "content")

		app := newTestApp(t)
		registerRoute(app, http.MethodPost, "/books", handler.Create, withLocals(locals))

		resp := performRequest(t, app, http.MethodPost, "/books", body, contentType)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "upload failed", envelope.Message)
	})

	t.Run("success", func(t *testing.T) {
		mockService := booksmocks.NewMockService(t)
		handler := NewBooks(mockService, stubUploader{url: "https://img/book.png"}, stubRecommendation{})

		mockService.EXPECT().
			Create(anyCtx(), mock.MatchedBy(func(req *books.BookRequest) bool {
				return req.ImageLink == "https://img/book.png" && req.Title == "Go"
			})).
			Return(books.Book{ID: 1, Title: "Go", ImageLink: "https://img/book.png"}, nil).
			Once()

		body, contentType := newMultipartBody(t, bookFields, "book.png", "content")

		app := newTestApp(t)
		registerRoute(app, http.MethodPost, "/books", handler.Create, withLocals(locals))

		resp := performRequest(t, app, http.MethodPost, "/books", body, contentType)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "book created", envelope.Message)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := booksmocks.NewMockService(t)
		handler := NewBooks(mockService, stubUploader{url: "https://img/book.png"}, stubRecommendation{})

		mockService.EXPECT().
			Create(anyCtx(), mock.MatchedBy(func(req *books.BookRequest) bool { return true })).
			Return(books.Book{}, errors.New("db down")).
			Once()

		body, contentType := newMultipartBody(t, bookFields, "book.png", "content")

		app := newTestApp(t)
		registerRoute(app, http.MethodPost, "/books", handler.Create, withLocals(locals))

		resp := performRequest(t, app, http.MethodPost, "/books", body, contentType)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestBooksUpdate(t *testing.T) {
	bookReq := &books.BookRequest{Title: "Go", Description: "d", Publisher: "p", Year: 2024, Stock: 1, CategoryID: 1}
	locals := map[string]any{"validatedBody": bookReq}

	t.Run("invalid id", func(t *testing.T) {
		handler, _ := testBooksHandler(t)

		app := newTestApp(t)
		registerRoute(app, http.MethodPatch, "/books/:id", handler.Update, withLocals(locals))

		resp := performRequest(t, app, http.MethodPatch, "/books/abc", nil, "")
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "invalid id", envelope.Message)
	})

	t.Run("file not found", func(t *testing.T) {
		handler, _ := testBooksHandler(t)

		body, contentType := newMultipartBody(t, nil, "", "")

		app := newTestApp(t)
		registerRoute(app, http.MethodPatch, "/books/:id", handler.Update, withLocals(locals))

		resp := performRequest(t, app, http.MethodPatch, "/books/1", body, contentType)
		require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "file not found", envelope.Message)
	})

	t.Run("invalid file format", func(t *testing.T) {
		handler, _ := testBooksHandler(t)

		body, contentType := newMultipartBody(t, nil, "book.gif", "content")

		app := newTestApp(t)
		registerRoute(app, http.MethodPatch, "/books/:id", handler.Update, withLocals(locals))

		resp := performRequest(t, app, http.MethodPatch, "/books/1", body, contentType)
		require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "invalid file format", envelope.Message)
	})

	t.Run("upload error", func(t *testing.T) {
		mockService := booksmocks.NewMockService(t)
		handler := NewBooks(mockService, stubUploader{err: errors.New("cloudinary down")}, stubRecommendation{})

		body, contentType := newMultipartBody(t, nil, "book.png", "content")

		app := newTestApp(t)
		registerRoute(app, http.MethodPatch, "/books/:id", handler.Update, withLocals(locals))

		resp := performRequest(t, app, http.MethodPatch, "/books/1", body, contentType)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "upload failed", envelope.Message)
	})

	t.Run("success", func(t *testing.T) {
		mockService := booksmocks.NewMockService(t)
		handler := NewBooks(mockService, stubUploader{url: "https://img/book.png"}, stubRecommendation{})

		mockService.EXPECT().
			Update(anyCtx(), mock.MatchedBy(func(req *books.BookRequest) bool {
				return req.ImageLink == "https://img/book.png"
			}), uint(1)).
			Return(books.Book{ID: 1, Title: "Go", ImageLink: "https://img/book.png"}, nil).
			Once()

		body, contentType := newMultipartBody(t, nil, "book.png", "content")

		app := newTestApp(t)
		registerRoute(app, http.MethodPatch, "/books/:id", handler.Update, withLocals(locals))

		resp := performRequest(t, app, http.MethodPatch, "/books/1", body, contentType)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "book updated", envelope.Message)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := booksmocks.NewMockService(t)
		handler := NewBooks(mockService, stubUploader{url: "https://img/book.png"}, stubRecommendation{})

		mockService.EXPECT().
			Update(anyCtx(), mock.MatchedBy(func(req *books.BookRequest) bool { return true }), uint(1)).
			Return(books.Book{}, errors.New("db down")).
			Once()

		body, contentType := newMultipartBody(t, nil, "book.png", "content")

		app := newTestApp(t)
		registerRoute(app, http.MethodPatch, "/books/:id", handler.Update, withLocals(locals))

		resp := performRequest(t, app, http.MethodPatch, "/books/1", body, contentType)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestBooksDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		handler, mockService := testBooksHandler(t)
		mockService.EXPECT().
			Delete(anyCtx(), uint(1)).
			Return(nil).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodDelete, "/books/:id", handler.Delete)

		resp := performRequest(t, app, http.MethodDelete, "/books/1", nil, "")
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "book deleted", envelope.Message)
	})

	t.Run("invalid id", func(t *testing.T) {
		handler, _ := testBooksHandler(t)

		app := newTestApp(t)
		registerRoute(app, http.MethodDelete, "/books/:id", handler.Delete)

		resp := performRequest(t, app, http.MethodDelete, "/books/abc", nil, "")
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("service error", func(t *testing.T) {
		handler, mockService := testBooksHandler(t)
		mockService.EXPECT().
			Delete(anyCtx(), uint(1)).
			Return(errors.New("not found")).
			Once()

		app := newTestApp(t)
		registerRoute(app, http.MethodDelete, "/books/:id", handler.Delete)

		resp := performRequest(t, app, http.MethodDelete, "/books/1", nil, "")
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestBooksGetBookRecommendation(t *testing.T) {
	recommReq := &ai.BookRecommendationRequest{Quantity: 2, Topic: "go"}
	locals := map[string]any{"validatedBody": recommReq}

	t.Run("success", func(t *testing.T) {
		mockService := booksmocks.NewMockService(t)
		content := `[{"number":1,"title":"Go in Action","summary":"s","publisher":"p","year":2015}]`
		handler := NewBooks(mockService, stubUploader{}, stubRecommendation{
			response: ai.PromptResponse{Choices: []ai.Choice{{Message: ai.ResponseMessage{Content: content}}}},
		})

		app := newTestApp(t)
		registerRoute(app, http.MethodPost, "/books/recommendations", handler.GetBookRecommendation, withLocals(locals))

		resp := performRequest(t, app, http.MethodPost, "/books/recommendations", nil, "")
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "book recommendations", envelope.Message)
	})

	t.Run("recommendation service error", func(t *testing.T) {
		mockService := booksmocks.NewMockService(t)
		handler := NewBooks(mockService, stubUploader{}, stubRecommendation{err: errors.New("ai down")})

		app := newTestApp(t)
		registerRoute(app, http.MethodPost, "/books/recommendations", handler.GetBookRecommendation, withLocals(locals))

		resp := performRequest(t, app, http.MethodPost, "/books/recommendations", nil, "")
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "get book recommendation failed. Please retry soon", envelope.Message)
	})

	t.Run("invalid recommendation format", func(t *testing.T) {
		mockService := booksmocks.NewMockService(t)
		handler := NewBooks(mockService, stubUploader{}, stubRecommendation{
			response: ai.PromptResponse{Choices: []ai.Choice{{Message: ai.ResponseMessage{Content: "not-json"}}}},
		})

		app := newTestApp(t)
		registerRoute(app, http.MethodPost, "/books/recommendations", handler.GetBookRecommendation, withLocals(locals))

		resp := performRequest(t, app, http.MethodPost, "/books/recommendations", nil, "")
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var envelope responseEnvelope
		readResponse(t, resp, &envelope)
		require.Equal(t, "get book recommendation failed. Invalid response", envelope.Message)
	})
}
