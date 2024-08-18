package repository

import (
	"aws-ses-local-go/domain"
	"aws-ses-local-go/internal/dao"
)

type MailRepository struct {
	MailDao dao.MailDao
}

func NewMailRepository(
	mailDao dao.MailDao,
) *MailRepository {
	return &MailRepository{
		MailDao: mailDao,
	}
}

func (r *MailRepository) Store(mail domain.Mail) error {
	err := r.MailDao.Store(mail)
	if err != nil {
		return err
	}
	return nil
}
