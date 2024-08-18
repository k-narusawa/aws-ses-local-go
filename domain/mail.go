package domain

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/mail"
	"strings"
	"time"

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
	CreatedAt           time.Time
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
		CreatedAt:           time.Now(),
	}
}

func FromRawEmailRequest(rawMessage string) (Mail, error) {
	decodedData, err := base64.StdEncoding.DecodeString(rawMessage)
	if err != nil {
		log.Printf("failed to decode base64: %v", err)
		return Mail{}, err
	}

	message, err := parseRawEmail(string(decodedData))
	if err != nil {
		log.Printf("failed to parse raw email: %v", err)
		return Mail{}, err
	}
	to := message.Header.Get("To")
	listUnsubscribeUrl := strings.Join(message.Header["List-Unsubscribe"], ",")
	listUnsubscribePost := strings.Join(message.Header["List-Unsubscribe-Post"], ",")
	body, err := getBody(message)
	if err != nil {
		return Mail{}, err
	}

	return Mail{
		MessageID:           generateMessageID(),
		From:                message.Header.Get("From"),
		To:                  &to,
		Subject:             message.Header.Get("Subject"),
		Text:                &body,
		ListUnsubscribeUrl:  &listUnsubscribeUrl,
		ListUnsubscribePost: &listUnsubscribePost,
		CreatedAt:           time.Now(),
	}, nil
}

func generateMessageID() string {
	seed := rand.Int63()
	rand.New(rand.NewSource(seed))
	messageId := fmt.Sprintf("ses-%d", rand.Intn(900000000)+100000000)
	return messageId
}

func parseRawEmail(rawEmail string) (*mail.Message, error) {
	reader := bytes.NewReader([]byte(rawEmail))
	message, err := mail.ReadMessage(reader)
	if err != nil {
		log.Printf("failed to parse raw email: %v", err)
		return nil, err
	}
	return message, nil
}

func getBody(message *mail.Message) (string, error) {
	bodyBytes, err := io.ReadAll(message.Body)
	if err != nil {
		log.Printf("failed to read body: %v", err)
		return "", err
	}
	return string(bodyBytes), nil
}
