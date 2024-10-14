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

func SendSimpleEmail(client *sesv2.Client) {
	input := &sesv2.SendEmailInput{
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Text: &types.Content{
						Data: aws.String("こんにちは"),
					},
				},
				Subject: &types.Content{
					Data: aws.String("Hello"),
				},
			},
		},
		Destination: &types.Destination{
			ToAddresses: []string{"to@example.com"},
		},
		FromEmailAddress: aws.String("from@example.com"),
	}

	result, err := client.SendEmail(context.TODO(), input)
	if err != nil {
		panic(err)
	}

	fmt.Println("SendSimpleEmail is OK. MessageID: ", *result.MessageId)
}

func SendRawEmail(client *sesv2.Client) {
	rawMsg, _ := formatRawEmailMessage()
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

func formatRawEmailMessage() ([]byte, error) {
	message := gomail.NewMessage()
	body := []byte("こんにちは")

	message.SetHeader("From", "from@example.com")
	message.SetHeader("To", "to@example.com")
	message.SetHeader("Subject", "こんにちは")
	message.SetBody("text/plain", string(body[:]))

	buf := new(bytes.Buffer)
	_, err := message.WriteTo(buf)
	if err != nil {
		return buf.Bytes(), errors.Wrap(err, "Failed to generate raw email notification message")
	}

	return buf.Bytes(), nil
}
