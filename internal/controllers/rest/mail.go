package rest

import (
	"aws-ses-local-go/usecase/query"
	"strconv"

	"github.com/labstack/echo/v4"
)

type IMailService interface {
	DeleteMail(mId string) error
	DeleteMails() error
}

type MailHandler struct {
	MailDtoQueryService query.IMailDtoQueryService
	CountQueryService   query.ICountQueryService
	MailService         IMailService
}

func NewMailHandler(
	e *echo.Echo,
	iMailDtoQueryService query.IMailDtoQueryService,
	iCountQueryService query.ICountQueryService,
	MailService IMailService,
) {
	handler := &MailHandler{
		MailDtoQueryService: iMailDtoQueryService,
		CountQueryService:   iCountQueryService,
		MailService:         MailService,
	}

	e.GET("/store", handler.GetMails)
	e.GET("/emails", handler.GetMails)
	e.DELETE("/emails/:message_id", handler.DeleteMail)
	e.DELETE("/emails", handler.DeleteMails)
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

	if mails == nil {
		mails = []query.MailDto{}
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

func (h *MailHandler) DeleteMail(c echo.Context) error {
	mId := c.Param("message_id")

	err := h.MailService.DeleteMail(mId)
	if err != nil {
		return c.JSON(500, err)
	}

	return c.NoContent(204)
}

func (h *MailHandler) DeleteMails(c echo.Context) error {
	err := h.MailService.DeleteMails()
	if err != nil {
		return c.JSON(500, err)
	}

	return c.NoContent(204)
}
