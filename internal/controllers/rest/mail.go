package rest

import (
	"aws-ses-local-go/usecase/query"
	"strconv"

	"github.com/labstack/echo/v4"
)

type MailHandler struct {
	MailDtoQueryService query.IMailDtoQueryService
	CountQueryService   query.ICountQueryService
}

func NewMailHandler(
	e *echo.Echo,
	iMailDtoQueryService query.IMailDtoQueryService,
	iCountQueryService query.ICountQueryService,
) {
	handler := &MailHandler{
		MailDtoQueryService: iMailDtoQueryService,
		CountQueryService:   iCountQueryService,
	}

	e.GET("/store", handler.GetMails)
	e.GET("/emails", handler.GetMails)
}

type MailResponse struct {
	Page      int             `json:"page"`
	Limit     int             `json:"limit"`
	Size      int             `json:"size"`
	TotalPage int             `json:"total_page"`
	TotalSize int             `json:"total_size"`
	Items     []query.MailDto `json:"items"`
}

func (h *MailHandler) GetMails(c echo.Context) error {
	qPage := c.QueryParam("page")
	qLimit := c.QueryParam("limit")
	to := c.QueryParam("to_address")

	if qPage == "" {
		qPage = "1"
	}

	if qLimit == "" {
		qLimit = "10"
	}

	orgSize, _ := strconv.Atoi(qPage)
	isize := orgSize - 1
	limit, _ := strconv.Atoi(qLimit)
	offset := isize * limit

	mails, err := h.MailDtoQueryService.FindByTo(&to, limit, offset)
	if err != nil {
		return c.JSON(500, err)
	}

	totalSize, err := h.CountQueryService.CountByTo(&to)
	if err != nil {
		return c.JSON(500, err)
	}

	resp := MailResponse{
		Page:      orgSize,
		Limit:     limit,
		Size:      len(mails),
		TotalPage: totalSize / limit,
		TotalSize: totalSize,
		Items:     mails,
	}

	return c.JSON(200, resp)
}
