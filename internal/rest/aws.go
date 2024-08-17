package rest

import (
	"aws-ses-local-go/usecase/aws"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
)

type AwsService interface {
	SendEmail(aws.SendEmailInput) (*aws.SendEmailOutput, error)
}

type AwsHandler struct {
	AwsService AwsService
}

func NewAwsHandler(e *echo.Echo, awsService AwsService) {
	handler := &AwsHandler{
		AwsService: awsService,
	}

	e.POST("", handler.SendEmail)
}

func (h *AwsHandler) SendEmail(c echo.Context) error {
	in := aws.SendEmailInput{
		Action:               c.QueryParam("Action"),
		Version:              c.QueryParam("Version"),
		ConfigurationSetName: c.QueryParam("ConfigurationSetName"),
		ToAddresses:          c.QueryParam("Destination.ToAddresses.member.1"),
		CcAddresses:          c.QueryParam("Destination.CcAddresses.member.1"),
		BccAddresses:         c.QueryParam("Destination.BccAddresses.member.1"),
		HtmlData:             c.QueryParam("Message.Body.Html.Data"),
		HtmlCharset:          c.QueryParam("Message.Body.Html.Charset"),
		TextData:             c.QueryParam("Message.Body.Text.Data"),
		TextCharset:          c.QueryParam("Message.Body.Text.Charset"),
		SubjectData:          c.QueryParam("Message.Subject.Data"),
		SubjectCharset:       c.QueryParam("Message.Subject.Charset"),
		ReplyToAddresses:     c.QueryParam("ReplyToAddresses.member.1"),
		ReturnPath:           c.QueryParam("ReturnPath"),
		ReturnPathArn:        c.QueryParam("ReturnPathArn"),
		Source:               c.QueryParam("Source"),
		SourceArn:            c.QueryParam("SourceArn"),
		Tags:                 c.QueryParam("Tags.member.1"),
		Destination:          c.QueryParam("Destination.member.1"),
		FromArn:              c.QueryParam("FromArn"),
		RawMessage:           c.QueryParam("RawMessage.Data"),
	}

	out, err := h.AwsService.SendEmail(in)
	if err != nil {
		return c.JSON(500, err)
	}

	log.Printf("MessageID: %s", out.MessageID)

	resp := fmt.Sprintf(`<SendEmailResponse xmlns="http://ses.amazonaws.com/doc/2010-12-01/"><SendEmailResult><MessageId>%s</MessageId></SendEmailResult></SendEmailResponse>`, out.MessageID)
	return c.XMLBlob(200, []byte(resp))
}
