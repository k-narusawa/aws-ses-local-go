package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String("http://localhost:8080"),
	})
	if err != nil {
		fmt.Println("Error creating session:")
		fmt.Println(err)
		return
	}

	svc := ses.New(sess)
	sendEmail(svc)
	sendRawEmail(svc)
}

func sendEmail(svc *ses.SES) {
	const (
		Sender    = "sender@example.com"
		Recipient = "recipient@example.com"
		Subject   = "Amazon SES Test (AWS SDK for Go)"
		HtmlBody  = "<h1>Amazon SES Test Email (AWS SDK for Go)</h1><p>This email was sent with " +
			"<a href='https://aws.amazon.com/ses/'>Amazon SES</a> using the " +
			"<a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.</p>"
		TextBody = "This email was sent with Amazon SES using the AWS SDK for Go."
		CharSet  = "UTF-8"
	)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(Recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(HtmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(Subject),
			},
		},
		Source: aws.String(Sender),
	}

	result, err := svc.SendEmail(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}

		return
	}

	fmt.Println("SendEmail is OK. MessageID: ", *result.MessageId)
}

func sendRawEmail(svc *ses.SES) {
	const (
		Sender    = "sender@example.com"
		Recipient = "recipient@example.com"
		Subject   = "Amazon SES Test (AWS SDK for Go)"
		HtmlBody  = "<h1>Amazon SES Test Email (AWS SDK for Go)</h1><p>This email was sent with " +
			"<a href='https://aws.amazon.com/ses/'>Amazon SES</a> using the " +
			"<a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.</p>"
		TextBody = "This email was sent with Amazon SES using the AWS SDK for Go."
		CharSet  = "UTF-8"
	)

	input := &ses.SendRawEmailInput{
		RawMessage: &ses.RawMessage{
			Data: []byte("From: " + Sender + "\n" + // Required
				"To: " + Recipient + "\n" + // Required
				"Subject: " + Subject + "\n" +
				"MIME-Version: 1.0\n" +
				"Content-type: Multipart/Mixed; boundary=\"NextPart\"\n\n" +
				"--NextPart\n" +
				"Content-Type: text/plain\n\n" +
				TextBody + "\n\n" +
				"--NextPart\n" +
				"Content-Type: text/html\n\n" +
				HtmlBody + "\n\n" +
				"--NextPart--"),
		},
	}

	result, err := svc.SendRawEmail(input)
	if err != nil {
		log.Println("Error sending email:")
		return
	}

	fmt.Println("SendRawEmail is OK. MessageID: ", *result.MessageId)
}
