package dao

import (
	"aws-ses-local-go/domain"

	"gorm.io/gorm"
)

type MailDao struct {
	db *gorm.DB
}

func NewMailDao(db *gorm.DB) *MailDao {
	return &MailDao{db: db}
}

func (d *MailDao) Store(mail domain.Mail) error {
	return d.db.Create(&mail).Error
}

func (d *MailDao) FindByMessageID(messageID string) (domain.Mail, error) {
	var mail domain.Mail

	err := d.db.Where("message_id = ?", messageID).First(&mail).Error
	if err != nil {
		return domain.Mail{}, err
	}

	return mail, err
}

func (d *MailDao) FindAll(limit, offset int) ([]domain.Mail, error) {
	var mails []domain.Mail
	err := d.db.Limit(limit).Offset(offset).Find(&mails).Error
	if err != nil {
		return nil, err
	}

	return mails, nil
}

func (d *MailDao) FindByTo(to *string, limit, offset int) ([]domain.Mail, error) {
	var mails []domain.Mail

	if to == nil || *to == "" {
		err := d.db.Limit(limit).Offset(offset).Find(&mails).Error
		if err != nil {
			return nil, err
		}

		return mails, nil
	}

	err := d.db.Where("`to` = ?", &to).Limit(limit).Offset(offset).Find(&mails).Error
	if err != nil {
		return nil, err
	}

	return mails, nil
}
