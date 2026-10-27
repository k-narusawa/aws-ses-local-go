// Package testutil はテスト用の共通ヘルパーを提供する。
package testutil

import (
	"fmt"
	"sync/atomic"
	"testing"

	"aws-ses-local-go/domain"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbSeq atomic.Int64

// NewTestDB はテストごとに独立したインメモリ SQLite を作成し、mails テーブルをマイグレーションする。
func NewTestDB(t testing.TB) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:testdb_%d?mode=memory&cache=shared", dbSeq.Add(1))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get sql.DB: %v", err)
	}
	// 共有キャッシュのインメモリ DB は最後の接続が閉じると消えるため、1 接続に固定する
	sqlDB.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = sqlDB.Close() })

	if err := db.AutoMigrate(&domain.Mail{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}
