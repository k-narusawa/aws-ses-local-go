package repository

import "aws-ses-local-go/domain"

var (
	mails = map[string]*domain.Mail{}
)

type MailRepository struct{}

func NewMailRepository() *MailRepository {
	return &MailRepository{}
}

func (r *MailRepository) Store(mail domain.Mail) error {
	mails[mail.MessageID] = &mail
	return nil
}

// <?xml version="1.0" encoding="UTF-8"?>
// <SendRawEmailResponse xmlns="http://ses.amazonaws.com/doc/2010-12-01/">
// 	<SendRawEmailResult>
// 		<MessageId>ses-899841303</MessageId>
// 	</SendRawEmailResult>
// </SendRawEmailResponse>
