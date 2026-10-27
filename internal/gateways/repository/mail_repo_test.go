package repository_test

import (
	"aws-ses-local-go/domain"
	"aws-ses-local-go/internal/gateways/dao"
	"aws-ses-local-go/internal/gateways/repository"
	"aws-ses-local-go/internal/testutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MailRepository が domain.IMailRepo を満たしていることをコンパイル時に保証する
var _ domain.IMailRepo = (*repository.MailRepository)(nil)

func TestMailRepository(t *testing.T) {
	d := dao.NewMailDao(testutil.NewTestDB(t))
	r := repository.NewMailRepository(*d)

	for _, id := range []string{"m1", "m2", "m3"} {
		require.NoError(t, r.Store(domain.Mail{MessageID: id, From: "f", Subject: "s", CreatedAt: time.Now()}))
	}

	require.NoError(t, r.Delete("m1"))
	_, err := d.FindByMessageID("m1")
	assert.Error(t, err)

	n, err := d.CountByTo(nil)
	require.NoError(t, err)
	assert.Equal(t, 2, n)

	require.NoError(t, r.DeleteAll())
	n, err = d.CountByTo(nil)
	require.NoError(t, err)
	assert.Equal(t, 0, n)
}

func TestMailRepository_Errors(t *testing.T) {
	db := testutil.NewTestDB(t)
	r := repository.NewMailRepository(*dao.NewMailDao(db))
	sqlDB, _ := db.DB()
	require.NoError(t, sqlDB.Close())

	assert.Error(t, r.Store(domain.Mail{MessageID: "x", From: "f", Subject: "s", CreatedAt: time.Now()}))
	assert.Error(t, r.Delete("x"))
	assert.Error(t, r.DeleteAll())
}
