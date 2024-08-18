package web

import (
	"aws-ses-local-go/usecase/query"

	"github.com/labstack/echo/v4"
)

type WebHandler struct {
	IMailDtoQueryService query.IMailDtoQueryService
}

func NewWebHandler(
	e *echo.Echo,
	iMailDtoQueryService query.IMailDtoQueryService,
) {
	handler := &WebHandler{
		IMailDtoQueryService: iMailDtoQueryService,
	}

	e.GET("/", handler.Index)
}

func (h *WebHandler) Index(c echo.Context) error {
	dummyData := []map[string]interface{}{
		{
			"id":      1,
			"message": "Hello Go",
		},
		{
			"id":      2,
			"message": "Hello again, Go",
		},
		{
			"id":      3,
			"message": "Hello once more, Go",
		},
	}
	return c.Render(200, "index.html", dummyData)
}
