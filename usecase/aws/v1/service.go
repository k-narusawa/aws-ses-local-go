package v1

import (
	"aws-ses-local-go/domain"
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

func (s *Service) SendEmail(in SendEmailInput) (*SendEmailOutput, error) {
	mail := domain.NewMail(
		in.Source,
		&in.ToAddresses,
		&in.CcAddresses,
		&in.BccAddresses,
		in.SubjectData,
		&in.TextData,
		&in.HtmlData,
		nil,
		nil,
	)

	err := s.MailRepo.Store(mail)
	if err != nil {
		return nil, err
	}

	return &SendEmailOutput{MessageID: mail.MessageID}, nil
}

type SendRawEmailInput struct {
	Version     string
	Source      string
	SourceArn   string
	Tags        string
	Destination string
	FromArn     string
	RawMessage  string
}

func (s *Service) SendRawEmail(in SendRawEmailInput) (*SendEmailOutput, error) {
	mail, err := domain.FromRawEmailRequest(in.RawMessage)
	if err != nil {
		return nil, err
	}

	err = s.MailRepo.Store(mail)
	if err != nil {
		return nil, err
	}

	return &SendEmailOutput{MessageID: mail.MessageID}, nil
}
