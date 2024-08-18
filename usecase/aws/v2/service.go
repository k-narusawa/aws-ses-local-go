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

	var text *string
	var html *string

	if in.Content.Simple.Body.Text != nil {
		text = &in.Content.Simple.Body.Text.Data
	} else if in.Content.Simple.Body.Html != nil {
		html = &in.Content.Simple.Body.Html.Data
	}

	mail := domain.NewMail(
		*in.FromEmailAddress,
		&to,
		&cc,
		&bcc,
		in.Content.Simple.Subject.Data,
		text,
		html,
		nil,
		nil,
	)

	err := s.MailRepo.Store(mail)
	if err != nil {
		return nil, err
	}

	return &SendEmailV2Output{MessageID: mail.MessageID}, nil
}

func (s *Service) SendRawEmail(in V2EmailOutboundEmailInput) (*SendEmailV2Output, error) {
	mail, err := domain.FromRawEmailRequest(in.Content.Raw.Data)
	if err != nil {
		return nil, err
	}

	err = s.MailRepo.Store(mail)
	if err != nil {
		return nil, err
	}

	return &SendEmailV2Output{MessageID: mail.MessageID}, nil
}
