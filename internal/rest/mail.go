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

type MailResponse struct {
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
	Size  int             `json:"size"`
	Items []query.MailDto `json:"items"`
}

func (h *MailHandler) GetMails(c echo.Context) error {
	qPage := c.QueryParam("page")
	qLimit := c.QueryParam("limit")
	to := c.QueryParam("to_address")

	if qPage == "" {
		qPage = "0"
	}

	if qLimit == "" {
		qLimit = "10"
	}

	isize, _ := strconv.Atoi(qPage)
	limit, _ := strconv.Atoi(qLimit)
	offset := isize * limit

	mails, err := h.MailDtoQueryService.FindByTo(&to, limit, offset)
	if err != nil {
		return c.JSON(500, err)
	}

	resp := MailResponse{
		Page:  isize,
		Limit: limit,
		Size:  len(mails),
		Items: mails,
	}

	return c.JSON(200, resp)
}
