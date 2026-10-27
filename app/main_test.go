package main

import (
	"net/http"
	"regexp"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterViews(t *testing.T) {
	e := echo.New()
	registerViews(e)

	get := func(path string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		return rec
	}

	rec := get("/")
	require.Equal(t, 200, rec.Code)
	assert.Contains(t, rec.Header().Get(echo.HeaderContentType), "text/html")
	assert.Contains(t, rec.Body.String(), `<div id="root">`)

	// index.html が参照している JS / CSS も埋め込まれていること
	refs := regexp.MustCompile(`(?:src|href)="(/[^"]+)"`).FindAllStringSubmatch(rec.Body.String(), -1)
	require.NotEmpty(t, refs)
	for _, ref := range refs {
		assert.Equal(t, 200, get(ref[1]).Code, ref[1])
	}

	rec = get("/favicon.ico")
	assert.Equal(t, 200, rec.Code)

	assert.Equal(t, 404, get("/not-found.txt").Code)
}

func TestHealthCheck(t *testing.T) {
	e := echo.New()
	e.GET("/health", healthCheck)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())
}
