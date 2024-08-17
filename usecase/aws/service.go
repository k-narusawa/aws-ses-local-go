package aws

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

type SendEmailInput struct {
	Action               string
	Version              string
	ConfigurationSetName string
	ToAddresses          string
	CcAddresses          string
	BccAddresses         string
	HtmlData             string
	HtmlCharset          string
	TextData             string
	TextCharset          string
	SubjectData          string
	SubjectCharset       string
	ReplyToAddresses     string
	ReturnPath           string
	ReturnPathArn        string
	Source               string
	SourceArn            string
	Tags                 string
	Destination          string
	FromArn              string
	RawMessage           string
}

type SendEmailOutput struct {
	MessageID string
}

func (s *Service) SendEmail(in SendEmailInput) (*SendEmailOutput, error) {
	if in.Action == "SendEmail" {
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

	mail := domain.FromRawEmailRequest(in.RawMessage)

	err := s.MailRepo.Store(mail)
	if err != nil {
		return nil, err
	}

	return &SendEmailOutput{MessageID: mail.MessageID}, nil
}
