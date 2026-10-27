package config_test

import (
	"aws-ses-local-go/config"
	"aws-ses-local-go/domain"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDBConnect_CreatesDBFileAndTable(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("LIMIT_RECORDS", "")

	db := config.DBConnect()
	sqlDB, err := db.DB()
	require.NoError(t, err)
	t.Cleanup(func() { _ = sqlDB.Close() })

	_, err = os.Stat("ses.db")
	require.NoError(t, err)
	assert.True(t, db.Migrator().HasTable(&domain.Mail{}))
}

func TestDBConnect_LimitRecords(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	t.Setenv("LIMIT_RECORDS", "2")

	db := config.DBConnect()
	sqlDB, err := db.DB()
	require.NoError(t, err)
	t.Cleanup(func() { _ = sqlDB.Close() })

	base := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 1; i <= 4; i++ {
		require.NoError(t, db.Create(&domain.Mail{
			MessageID: fmt.Sprintf("m%d", i),
			From:      "f",
			Subject:   "s",
			CreatedAt: base.Add(time.Duration(i) * time.Minute),
		}).Error)
	}

	var ids []string
	require.NoError(t, db.Model(&domain.Mail{}).Order("created_at desc").Pluck("message_id", &ids).Error)
	// 新しい 2 件だけが残る
	assert.Equal(t, []string{"m4", "m3"}, ids)
	assert.FileExists(t, filepath.Join(dir, "ses.db"))
}
