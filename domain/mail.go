package domain

import (
	"fmt"

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

func generateMessageID() string {
	seed := rand.Int63()
	rand.New(rand.NewSource(seed))
	messageId := fmt.Sprintf("ses-%d", rand.Intn(900000000)+100000000)
	return messageId
}
