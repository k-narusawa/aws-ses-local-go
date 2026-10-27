package middleware_test

import (
	"aws-ses-local-go/domain"
	"aws-ses-local-go/internal/gateways/middleware"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleError(t *testing.T) {
	log.SetOutput(io.Discard)

	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   string
	}{
		{"user not found", domain.ErrUserNotFound, 404, `{"code":"404-001","message":"user not found."}`},
		{"todo not found", domain.ErrToDoNotFound, 404, `{"code":"404-002","message":"todo not found."}`},
		{"wrap されたエラーも unwrap して判定する", fmt.Errorf("outer: %w", fmt.Errorf("inner: %w", domain.ErrToDoNotFound)), 404, `{"code":"404-002","message":"todo not found."}`},
		{"その他は 500", errors.New("boom"), 500, `{"code":"500-001","message":"boom"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			rec := httptest.NewRecorder()
			c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec)

			require.NoError(t, middleware.HandleError(c, tt.err))
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}
