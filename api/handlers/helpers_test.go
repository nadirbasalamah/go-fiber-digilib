package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newTestApp(t *testing.T) *fiber.App {
	t.Helper()

	return fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, _ error) error {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "internal error"})
		},
	})
}

type testOption func(c fiber.Ctx)

func anyCtx() any {
	return mock.MatchedBy(func(_ context.Context) bool { return true })
}

func withLocals(locals map[string]any) testOption {
	return func(c fiber.Ctx) {
		for k, v := range locals {
			c.Locals(k, v)
		}
	}
}

func registerRoute(app *fiber.App, method, path string, handler fiber.Handler, options ...testOption) {
	var chain []any

	for _, opt := range options {
		opt := opt
		chain = append(chain, func(c fiber.Ctx) error {
			opt(c)

			return c.Next()
		})
	}

	chain = append(chain, handler)

	app.Add([]string{method}, path, chain[0], chain[1:]...)
}

func performRequest(t *testing.T, app *fiber.App, method, path string, body []byte, contentType string, options ...testOption) *http.Response {
	t.Helper()

	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}

	req := httptest.NewRequest(method, path, reader)

	if contentType != "" {
		req.Header.Set(fiber.HeaderContentType, contentType)
	}

	resp, err := app.Test(req)
	require.NoError(t, err)

	return resp
}

func performJSON(t *testing.T, method, path string, body []byte, handler fiber.Handler, options ...testOption) *http.Response {
	t.Helper()

	app := newTestApp(t)
	registerRoute(app, method, path, handler, options...)

	return performRequest(t, app, method, path, body, fiber.MIMEApplicationJSON)
}
func readResponse(t *testing.T, resp *http.Response, out any) {
	t.Helper()

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	require.NoError(t, json.Unmarshal(body, out), string(body))
}

func jsonBody(t *testing.T, v any) []byte {
	t.Helper()

	b, err := json.Marshal(v)
	require.NoError(t, err)

	return b
}

func newMultipartBody(t *testing.T, fields map[string]string, fileName, fileContent string) ([]byte, string) {
	t.Helper()

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	for k, v := range fields {
		require.NoError(t, writer.WriteField(k, v))
	}

	if fileName != "" {
		part, err := writer.CreateFormFile("file", fileName)
		require.NoError(t, err)
		_, err = io.Copy(part, strings.NewReader(fileContent))
		require.NoError(t, err)
	}

	require.NoError(t, writer.Close())

	return buf.Bytes(), writer.FormDataContentType()
}

type responseEnvelope struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func jsonUnmarshal(data []byte, out any) error {
	return json.Unmarshal(data, out)
}

func readerFrom(body []byte) io.Reader {
	if body == nil {
		return nil
	}

	return bytes.NewReader(body)
}
