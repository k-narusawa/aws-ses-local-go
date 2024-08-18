package config

import (
	"aws-ses-local-go/internal/dao"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func DBConnect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("ses.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&dao.MailRecord{})

	return db
}
