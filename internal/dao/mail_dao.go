package dao

import "aws-ses-local-go/domain"

var (
	mails = map[string]*domain.Mail{}
)

type MailDao struct{}

func NewMailDao() *MailDao {
	return &MailDao{}
}

func (d *MailDao) Store(mail domain.Mail) error {
	mails[mail.MessageID] = &mail
	return nil
}
