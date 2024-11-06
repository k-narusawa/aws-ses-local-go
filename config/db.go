package config

import (
	"aws-ses-local-go/domain"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func DBConnect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("ses.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&domain.Mail{})

	limit := os.Getenv("LIMIT_RECORDS")
	if limit != "" {
		log.Println("Limiting records to " + limit)
		triggerSQL := `
			DROP TRIGGER IF EXISTS limit_records_after_insert;
			CREATE TRIGGER limit_records_after_insert
			AFTER INSERT ON mails
			BEGIN
				DELETE FROM mails
				WHERE message_id IN (
					SELECT message_id FROM mails
					ORDER BY created_at DESC
					LIMIT 1 OFFSET ` + limit + `
				);
			END;
			`
		db.Exec(triggerSQL)
	}

	return db
}
