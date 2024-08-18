package v2

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/aws/aws-sdk-go/aws"
)

func SendSimpleEmail(client *sesv2.Client) {
	input := &sesv2.SendEmailInput{
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Text: &types.Content{
						Data: aws.String("Hello, World!"),
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

	fmt.Println("SendEmail is OK. MessageID: ", *result.MessageId)
}
