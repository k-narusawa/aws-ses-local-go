package domain

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net/mail"
	"strings"

	"math/rand"
)

type Mail struct {
	MessageID           string
	From                string
	To                  *string
	Cc                  *string
	Bcc                 *string
	Subject             string
	Text                *string
	Html                *string
	ListUnsubscribePost *string
	ListUnsubscribeUrl  *string
}

func NewMail(
	from string,
	to *string,
	cc *string,
	bcc *string,
	subject string,
	text *string,
	html *string,
	listUnsubscribePost *string,
	listUnsubscribeUrl *string,
) Mail {
	return Mail{
		MessageID:           generateMessageID(),
		From:                from,
		To:                  to,
		Cc:                  cc,
		Bcc:                 bcc,
		Subject:             subject,
		Text:                text,
		Html:                html,
		ListUnsubscribePost: listUnsubscribePost,
		ListUnsubscribeUrl:  listUnsubscribeUrl,
	}
}

func FromRawEmailRequest(rawMessage string) Mail {
	decodedData, err := base64.StdEncoding.DecodeString(rawMessage)
	if err != nil {
		log.Fatalf("failed to decode base64: %v", err)
	}

	message := parseRawEmail(string(decodedData))
	to := message.Header.Get("To")
	listUnsubscribeUrl := strings.Join(message.Header["List-Unsubscribe"], ",")
	listUnsubscribePost := strings.Join(message.Header["List-Unsubscribe-Post"], ",")
	body := getBody(message)

	return Mail{
		MessageID:           generateMessageID(),
		From:                message.Header.Get("From"),
		To:                  &to,
		Subject:             message.Header.Get("Subject"),
		Text:                &body,
		ListUnsubscribeUrl:  &listUnsubscribeUrl,
		ListUnsubscribePost: &listUnsubscribePost,
	}
}

func generateMessageID() string {
	seed := rand.Int63()
	rand.New(rand.NewSource(seed))
	messageId := fmt.Sprintf("ses-%d", rand.Intn(900000000)+100000000)
	return messageId
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
