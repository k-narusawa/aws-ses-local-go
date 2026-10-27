package web_test

import (
	"aws-ses-local-go/internal/controllers/web"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type stubRenderer struct {
	name string
	data interface{}
}

func (r *stubRenderer) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	r.name, r.data = name, data
	_, err := fmt.Fprint(w, "rendered")
	return err
}

func TestWebHandler_Index(t *testing.T) {
	e := echo.New()
	r := &stubRenderer{}
	e.Renderer = r
	web.NewWebHandler(e, nil)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "rendered", rec.Body.String())
	assert.Equal(t, "index.html", r.name)
	assert.Len(t, r.data, 3)
}
