package query

import (
	"aws-ses-local-go/internal/gateways/dao"
	"aws-ses-local-go/usecase/query"
)

type MailDtoQueryService struct {
	MailDao dao.MailDao
}

func NewMailDtoQueryService(
	mailDao dao.MailDao,
) *MailDtoQueryService {
	return &MailDtoQueryService{
		MailDao: mailDao,
	}
}

func (s *MailDtoQueryService) FindByMessageID(messageID string) (query.MailDto, error) {
	mail, err := s.MailDao.FindByMessageID(messageID)
	if err != nil {
		return query.MailDto{}, err
	}

	return query.MailDto{
		MessageID: mail.MessageID,
		From:      mail.From,
		Destination: query.Destination{
			To:  mail.To,
			Cc:  mail.Cc,
			Bcc: mail.Bcc,
		},
		Subject: mail.Subject,
		Body: query.Body{
			Text: mail.Text,
			Html: mail.Html,
		},
		ListUnsubscribePost: mail.ListUnsubscribePost,
		ListUnsubscribeUrl:  mail.ListUnsubscribeUrl,
		CreatedAt:           mail.CreatedAt,
	}, nil
}

func (s *MailDtoQueryService) FindAll(limit, offset int) ([]query.MailDto, error) {
	mails, err := s.MailDao.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	var mailDtos []query.MailDto
	for _, mail := range mails {
		mailDtos = append(mailDtos, query.MailDto{
			MessageID: mail.MessageID,
			From:      mail.From,
			Destination: query.Destination{
				To:  mail.To,
				Cc:  mail.Cc,
				Bcc: mail.Bcc,
			},
			Subject: mail.Subject,
			Body: query.Body{
				Text: mail.Text,
				Html: mail.Html,
			},
			ListUnsubscribePost: mail.ListUnsubscribePost,
			ListUnsubscribeUrl:  mail.ListUnsubscribeUrl,
			CreatedAt:           mail.CreatedAt,
		})
	}

	return mailDtos, nil
}

func (s *MailDtoQueryService) FindByTo(to *string, limit, offset int) ([]query.MailDto, error) {
	mails, err := s.MailDao.FindByTo(to, limit, offset)
	if err != nil {
		return nil, err
	}

	var mailDtos []query.MailDto
	for _, mail := range mails {
		mailDtos = append(mailDtos, query.MailDto{
			MessageID: mail.MessageID,
			From:      mail.From,
			Destination: query.Destination{
				To:  mail.To,
				Cc:  mail.Cc,
				Bcc: mail.Bcc,
			},
			Subject: mail.Subject,
			Body: query.Body{
				Text: mail.Text,
				Html: mail.Html,
			},
			ListUnsubscribePost: mail.ListUnsubscribePost,
			ListUnsubscribeUrl:  mail.ListUnsubscribeUrl,
			CreatedAt:           mail.CreatedAt,
		})
	}

	return mailDtos, nil
}
