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
