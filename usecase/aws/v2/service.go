package v2

import (
	"aws-ses-local-go/domain"
	"strings"
)

type Service struct {
	MailRepo domain.IMailRepo
}

func NewService(
	mailRepo domain.IMailRepo,
) *Service {
	return &Service{
		MailRepo: mailRepo,
	}
}

func (s *Service) SendSimpleEmail(in V2EmailOutboundEmailInput) (*SendEmailV2Output, error) {
	to := strings.Join(in.Destination.ToAddresses, ",")
	cc := strings.Join(in.Destination.CcAddresses, ",")
	bcc := strings.Join(in.Destination.BccAddresses, ",")
	mail := domain.NewMail(
		*in.FromEmailAddress,
		&to,
		&cc,
		&bcc,
		in.Content.Simple.Subject.Data,
		&in.Content.Simple.Body.Text.Data,
		&in.Content.Simple.Body.Html.Data,
		nil,
		nil,
	)

	err := s.MailRepo.Store(mail)
	if err != nil {
		return nil, err
	}

	return &SendEmailV2Output{MessageID: mail.MessageID}, nil
}
