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
	Subject            string
	Body               string
	From               string
	To                 string
	ListUnsubscribeUrl *string
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
				Headers: []types.MessageHeader{
					{
						Name:  aws.String("List-Unsubscribe"),
						Value: aws.String(fmt.Sprintf("<%s>", *in.ListUnsubscribeUrl)),
					},
					{
						Name:  aws.String("List-Unsubscribe-Post"),
						Value: aws.String("List-Unsubscribe=One-Click"),
					},
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
	Subject            string
	Body               string
	From               string
	To                 string
	ListUnsubscribeUrl *string
}

func SendRawEmail(client *sesv2.Client, in SendRawEmailInput) {
	rawMsg, _ := formatRawEmailMessage(in.Subject, in.Body, in.From, in.To, in.ListUnsubscribeUrl)
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

func formatRawEmailMessage(subject, body, from, to string, listUnsubscribeUrl *string) ([]byte, error) {
	message := gomail.NewMessage()
	messageBody := []byte(body)

	message.SetHeader("From", from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	if listUnsubscribeUrl != nil {
		message.SetHeader("List-Unsubscribe", *listUnsubscribeUrl)
		message.SetHeader("List-Unsubscribe-Post", "List-Unsubscribe=One-Click")
	}
	message.SetBody("text/plain", string(messageBody[:]))

	buf := new(bytes.Buffer)
	_, err := message.WriteTo(buf)
	if err != nil {
		return buf.Bytes(), errors.Wrap(err, "Failed to generate raw email notification message")
	}

	return buf.Bytes(), nil
}
