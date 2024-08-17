package aws

import (
	"aws-ses-local-go/domain"
	"bufio"
	"bytes"
	"encoding/base64"
	"log"
	"net/mail"
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
	log.Printf("Action: %s, RawMessage: %s", in.Action, in.RawMessage)
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

		return &SendEmailOutput{
			MessageID: mail.MessageID,
		}, nil
	}

	decodedData, err := base64.StdEncoding.DecodeString(in.RawMessage)

	if err != nil {
		log.Fatalf("failed to decode base64: %v", err)
	}

	message := parseRawEmail(string(decodedData))
	to := message.Header.Get("To")
	listUnsubscribeUrl := strings.Join(message.Header["List-Unsubscribe"], ",")
	listUnsubscribePost := strings.Join(message.Header["List-Unsubscribe-Post"], ",")

	body := getBody(message)
	mail := domain.NewMail(
		message.Header.Get("From"),
		&to,
		nil,
		nil,
		message.Header.Get("Subject"),
		&body,
		nil,
		&listUnsubscribeUrl,
		&listUnsubscribePost,
	)

	err = s.MailRepo.Store(mail)
	if err != nil {
		return nil, err
	}

	return &SendEmailOutput{
		MessageID: mail.MessageID,
	}, nil

}

func parseRawEmail(rawEmail string) *mail.Message {
	reader := bufio.NewReader(bytes.NewReader([]byte(rawEmail)))
	message, err := mail.ReadMessage(reader)
	if err != nil {
		log.Fatalf("failed to parse raw email: %v", err)
	}
	return message
}

func getBody(message *mail.Message) string {
	body, err := bufio.NewReader(message.Body).ReadString('\n')
	if err != nil {
		log.Fatalf("failed to read email body: %v", err)
	}
	return body
}
