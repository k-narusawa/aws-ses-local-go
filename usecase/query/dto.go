package query

import "time"

type MailDto struct {
	MessageID           string      `json:"message_id"`
	From                string      `json:"from"`
	Destination         Destination `json:"destination"`
	Subject             string      `json:"subject"`
	Body                Body        `json:"body"`
	ListUnsubscribePost *string     `json:"list_unsubscribe_post"`
	ListUnsubscribeUrl  *string     `json:"list_unsubscribe_url"`
	CreatedAt           time.Time   `json:"created_at"`
}

type Destination struct {
	To  *string `json:"to"`
	Cc  *string `json:"cc"`
	Bcc *string `json:"bcc"`
}

type Body struct {
	Text *string `json:"text"`
	Html *string `json:"html"`
}
