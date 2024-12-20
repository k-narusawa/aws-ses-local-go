package main

import (
	v1 "aws-ses-local-go/test/v1"
	v2 "aws-ses-local-go/test/v2"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

func main() {
	TestV1()
	TestV2()
}

func TestV1() {
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
	v1.SendEmail(svc)
	v1.SendRawEmail(svc)
}

func TestV2() {
	svc := sesv2.New(sesv2.Options{
		Region:       "us-west-2",
		BaseEndpoint: aws.String("http://localhost:8080"),
	})

	sin := v2.SendSimpleEmailInput{
		Subject:            "Hello from SESv2 SendSimpleEmail",
		Body:               "Hello from SESv2 SendSimpleEmail",
		From:               "from@example.com",
		To:                 "to@example.com",
		ListUnsubscribeUrl: aws.String("http://localhost:8080/test/unsubscribe"),
	}
	v2.SendSimpleEmail(svc, sin)

	rin := v2.SendRawEmailInput{
		Subject:            "Hello from SESv2 SendRawEmail",
		Body:               "Hello from SESv2 SendRawEmail",
		From:               "from@example.com",
		To:                 "to@example.com",
		ListUnsubscribeUrl: aws.String("http://localhost:8080/test/unsubscribe"),
	}
	v2.SendRawEmail(svc, rin)

	rin_2 := v2.SendRawEmailInput{
		Subject: "日本語のメール",
		Body:    "日本語のメール",
		From:    "from@example.com",
		To:      "to@example.com",
	}
	v2.SendRawEmail(svc, rin_2)
}
