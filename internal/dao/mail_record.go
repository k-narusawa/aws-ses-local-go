package dao

type MailRecord struct {
	MessageID           string  `gorm:"primaryKey"`
	From                string  `gorm:"not null"`
	To                  *string `gorm:"default:null"`
	Cc                  *string `gorm:"default:null"`
	Bcc                 *string `gorm:"default:null"`
	Subject             string  `gorm:"not null"`
	Text                *string `gorm:"default:null"`
	Html                *string `gorm:"default:null"`
	ListUnsubscribePost *string `gorm:"default:null"`
	ListUnsubscribeUrl  *string `gorm:"default:null"`
	CreatedAt           string  `gorm:"not null"`
}

func (MailRecord) TableName() string {
	return "mails"
}
