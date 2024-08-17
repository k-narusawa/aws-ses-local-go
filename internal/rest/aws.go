package rest

import (
	"aws-ses-local-go/usecase/aws"
	"fmt"

	"github.com/labstack/echo/v4"
)

type AwsService interface {
	SendEmail(aws.SendEmailInput) (*aws.SendEmailOutput, error)
	SendRawEmail(aws.SendRawEmailInput) (*aws.SendEmailOutput, error)
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
	action := c.FormValue("Action")

	if action == "SendEmail" {
		in := aws.SendEmailInput{
			Version:              c.FormValue("Version"),
			ConfigurationSetName: c.FormValue("ConfigurationSetName"),
			ToAddresses:          c.FormValue("Destination.ToAddresses.member.1"),
			CcAddresses:          c.FormValue("Destination.CcAddresses.member.1"),
			BccAddresses:         c.FormValue("Destination.BccAddresses.member.1"),
			HtmlData:             c.FormValue("Message.Body.Html.Data"),
			HtmlCharset:          c.FormValue("Message.Body.Html.Charset"),
			TextData:             c.FormValue("Message.Body.Text.Data"),
			TextCharset:          c.FormValue("Message.Body.Text.Charset"),
			SubjectData:          c.FormValue("Message.Subject.Data"),
			SubjectCharset:       c.FormValue("Message.Subject.Charset"),
			ReplyToAddresses:     c.FormValue("ReplyToAddresses.member.1"),
			ReturnPath:           c.FormValue("ReturnPath"),
			ReturnPathArn:        c.FormValue("ReturnPathArn"),
			Source:               c.FormValue("Source"),
			SourceArn:            c.FormValue("SourceArn"),
			Tags:                 c.FormValue("Tags.member.1"),
			Destination:          c.FormValue("Destination.member.1"),
			FromArn:              c.FormValue("FromArn"),
		}

		out, err := h.AwsService.SendEmail(in)
		if err != nil {
			return c.JSON(500, err)
		}

		resp := fmt.Sprintf(`<SendEmailResponse xmlns="http://ses.amazonaws.com/doc/2010-12-01/"><SendEmailResult><MessageId>%s</MessageId></SendEmailResult></SendEmailResponse>`, out.MessageID)
		return c.XMLBlob(200, []byte(resp))
	} else if action == "SendRawEmail" {
		rawMessage := c.FormValue("RawMessage.Data")
		in := aws.SendRawEmailInput{
			Version:     c.FormValue("Version"),
			Source:      c.FormValue("Source"),
			SourceArn:   c.FormValue("SourceArn"),
			Tags:        c.FormValue("Tags.member.1"),
			Destination: c.FormValue("Destination.member.1"),
			FromArn:     c.FormValue("FromArn"),
			RawMessage:  rawMessage,
		}

		out, err := h.AwsService.SendRawEmail(in)
		if err != nil {
			return c.JSON(500, err)
		}

		resp := fmt.Sprintf(`<SendRawEmailResponse xmlns="http://ses.amazonaws.com/doc/2010-12-01/"><SendRawEmailResult><MessageId>%s</MessageId></SendRawEmailResult></SendRawEmailResponse>`, out.MessageID)
		return c.XMLBlob(200, []byte(resp))
	}

	return c.JSON(400, "Invalid action")
}
