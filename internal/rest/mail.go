package rest

import (
	"aws-ses-local-go/usecase/query"
	"strconv"

	"github.com/labstack/echo/v4"
)

type MailHandler struct {
	MailDtoQueryService query.IMailDtoQueryService
}

func NewMailHandler(
	e *echo.Echo,
	iMailDtoQueryService query.IMailDtoQueryService,
) {
	handler := &MailHandler{
		MailDtoQueryService: iMailDtoQueryService,
	}

	e.GET("/store", handler.GetMails)
	e.GET("/emails", handler.GetMails)
}

func (h *MailHandler) GetMails(c echo.Context) error {
	page := c.QueryParam("page")
	size := c.QueryParam("size")
	to := c.QueryParam("to_address")

	if page == "" {
		page = "0"
	}

	if size == "" {
		size = "10"
	}

	limit, _ := strconv.Atoi(size)
	isize, _ := strconv.Atoi(page)
	offset := isize * limit

	mails, err := h.MailDtoQueryService.FindByTo(&to, limit, offset)
	if err != nil {
		return c.JSON(500, err)
	}

	return c.JSON(200, mails)
}
