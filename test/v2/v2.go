package v2

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
)

type SendSimpleEmailInput struct {
	Subject string
	Body    string
	From    string
	To      string
}

func SendSimpleEmail(client *sesv2.Client, in SendSimpleEmailInput) {
	input := &sesv2.SendEmailInput{
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Text: &types.Content{
						Data: aws.String(in.Subject),
					},
				},
				Subject: &types.Content{
					Data: aws.String(in.Subject),
				},
			},
		},
		Destination: &types.Destination{
			ToAddresses: []string{in.To},
		},
		FromEmailAddress: aws.String(in.From),
	}

	result, err := client.SendEmail(context.TODO(), input)
	if err != nil {
		panic(err)
	}

	fmt.Println("SendSimpleEmail is OK. MessageID: ", *result.MessageId)
}

type SendRawEmailInput struct {
	Subject string
	Body    string
	From    string
	To      string
}

func SendRawEmail(client *sesv2.Client, in SendRawEmailInput) {
	rawMsg, _ := formatRawEmailMessage(in.Subject, in.Body, in.From, in.To)
	input := &sesv2.SendEmailInput{
		Content: &types.EmailContent{
			Raw: &types.RawMessage{
				Data: []byte(rawMsg),
			},
		},
	}

	result, err := client.SendEmail(context.TODO(), input)
	if err != nil {
		panic(err)
	}

	fmt.Println("SendRawEmail is OK. MessageID: ", *result.MessageId)
}

func formatRawEmailMessage(subject, body, from, to string) ([]byte, error) {
	message := gomail.NewMessage()
	messageBody := []byte(body)

	message.SetHeader("From", from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", string(messageBody[:]))

	buf := new(bytes.Buffer)
	_, err := message.WriteTo(buf)
	if err != nil {
		return buf.Bytes(), errors.Wrap(err, "Failed to generate raw email notification message")
	}

	return buf.Bytes(), nil
}
