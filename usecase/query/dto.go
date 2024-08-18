package query

import "time"

type MailDto struct {
	MessageID           string    `json:"message_id"`
	From                string    `json:"from"`
	To                  *string   `json:"to"`
	Cc                  *string   `json:"cc"`
	Bcc                 *string   `json:"bcc"`
	Subject             string    `json:"subject"`
	Text                *string   `json:"text"`
	Html                *string   `json:"html"`
	ListUnsubscribePost *string   `json:"list_unsubscribe_post"`
	ListUnsubscribeUrl  *string   `json:"list_unsubscribe_url"`
	CreatedAt           time.Time `json:"created_at"`
}
