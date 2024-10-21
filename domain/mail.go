package domain

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/quotedprintable"
	"net/mail"
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
	_, params, err := mime.ParseMediaType(contentType)
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

	subject := message.Header.Get("Subject")
	decodedSubject, err := decodeHeader(subject)
	if err != nil {
		log.Printf("Failed to decode subject: %v", err)
	}

	contentTransferEncoding := message.Header.Get("Content-Transfer-Encoding")
	if contentTransferEncoding == "" {
		contentTransferEncoding = "7bit"
	}

	body, err := decodeBody(charset, message.Body, contentTransferEncoding)
	if err != nil {
		log.Printf("failed to decode body: %v", err)
		return Mail{}, err
	}

	to := message.Header.Get("To")
	listUnsubscribeUrl := message.Header.Get("List-Unsubscribe")
	listUnsubscribeUrl = strings.Trim(listUnsubscribeUrl, "<>")
	listUnsubscribePost := message.Header.Get("List-Unsubscribe-Post")

	return Mail{
		MessageID:           generateMessageID(),
		From:                message.Header.Get("From"),
		To:                  &to,
		Subject:             decodedSubject,
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

func decodeHeader(encoded string) (string, error) {
	decodedHeader, err := (&mime.WordDecoder{}).DecodeHeader(encoded)
	if err != nil {
		return "", err
	}
	return decodedHeader, nil
}

func decodeBody(charset string, body io.Reader, contentTransferEncoding string) (string, error) {
	var reader io.Reader
	var err error

	// Content-Transfer-Encodingに基づいてデコード処理を追加
	switch strings.ToLower(contentTransferEncoding) {
	case "base64":
		reader = base64.NewDecoder(base64.StdEncoding, body)
	case "quoted-printable":
		reader = quotedprintable.NewReader(body)
	default:
		reader = body // デフォルトは何もしない
	}

	// charsetに基づく変換
	reader, err = charsetlib.NewReader(reader, charset)
	if err != nil {
		return "", err
	}

	decodedBytes, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(decodedBytes), nil
}
