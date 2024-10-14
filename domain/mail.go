package domain

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime"
	"net/mail"
	"os"
	"strings"
	"time"

	"math/rand"

	charsetlib "golang.org/x/net/html/charset"
)

type Mail struct {
	MessageID           string    `gorm:"primaryKey"`
	From                string    `gorm:"not null"`
	To                  *string   `gorm:"default:null"`
	Cc                  *string   `gorm:"default:null"`
	Bcc                 *string   `gorm:"default:null"`
	Subject             string    `gorm:"not null"`
	Text                *string   `gorm:"default:null"`
	Html                *string   `gorm:"default:null"`
	ListUnsubscribePost *string   `gorm:"default:null"`
	ListUnsubscribeUrl  *string   `gorm:"default:null"`
	CreatedAt           time.Time `gorm:"not null"`
}

func (Mail) TableName() string {
	return "mails"
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

	contentType := message.Header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		log.Printf("Error parsing Content-Type: %v\n", err)
		return Mail{}, err
	}

	var charset string
	if cs, ok := params["charset"]; ok {
		charset = cs
	} else {
		charset = "utf-8" // デフォルト
	}

	if strings.HasPrefix(mediaType, "text/") {
		body, err := decodeBody(charset, msg.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error decoding body: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Body: %s\n", body)
	} else {
		fmt.Fprintf(os.Stderr, "Non-text content type: %s\n", contentType)
	}

	to := message.Header.Get("To")
	subject := message.Header.Get("Subject")
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
		Subject:             subject,
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

func decodeHeader(encoded string) (string, error) {
	decodedHeader, err := (&mime.WordDecoder{}).DecodeHeader(encoded)
	if err != nil {
		return "", err
	}
	return decodedHeader, nil
}

func decodeBody(charset string, body io.Reader) (string, error) {
	reader, err := charsetlib.NewReader(body, charset)
	if err != nil {
		return "", err
	}
	decodedBytes, err := io.ReadAll(reader)
	if err != nil {
		return ""çç, err
	}
	return string(decodedBytes), nil
}
