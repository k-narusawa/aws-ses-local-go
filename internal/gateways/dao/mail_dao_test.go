package dao_test

import (
	"aws-ses-local-go/domain"
	"aws-ses-local-go/internal/gateways/dao"
	"aws-ses-local-go/internal/testutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func ptr(s string) *string { return &s }

func newMail(id, to string, createdAt time.Time) domain.Mail {
	return domain.Mail{
		MessageID: id,
		From:      "from@example.com",
		To:        ptr(to),
		Subject:   "subject " + id,
		Text:      ptr("text " + id),
		CreatedAt: createdAt,
	}
}

// 3 件のメールを作成日時がずれる形で登録する。新しい順は m3, m2, m1。
func seed(t *testing.T, d *dao.MailDao) {
	t.Helper()
	base := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	require.NoError(t, d.Store(newMail("m1", "a@example.com", base)))
	require.NoError(t, d.Store(newMail("m2", "b@example.com", base.Add(time.Minute))))
	require.NoError(t, d.Store(newMail("m3", "a@example.com", base.Add(2*time.Minute))))
}

func ids(mails []domain.Mail) []string {
	out := make([]string, 0, len(mails))
	for _, m := range mails {
		out = append(out, m.MessageID)
	}
	return out
}

func TestMailDao_StoreAndFindByMessageID(t *testing.T) {
	d := dao.NewMailDao(testutil.NewTestDB(t))

	created := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)
	in := domain.Mail{
		MessageID:           "ses-123456789",
		From:                "from@example.com",
		To:                  ptr("to@example.com"),
		Cc:                  ptr("cc@example.com"),
		Bcc:                 ptr("bcc@example.com"),
		Subject:             "日本語の件名",
		Text:                ptr("text"),
		Html:                ptr("<p>html</p>"),
		ListUnsubscribePost: ptr("List-Unsubscribe=One-Click"),
		ListUnsubscribeUrl:  ptr("https://example.com/u"),
		CreatedAt:           created,
	}
	require.NoError(t, d.Store(in))

	got, err := d.FindByMessageID("ses-123456789")
	require.NoError(t, err)
	assert.Equal(t, in.MessageID, got.MessageID)
	assert.Equal(t, in.From, got.From)
	assert.Equal(t, *in.To, *got.To)
	assert.Equal(t, *in.Cc, *got.Cc)
	assert.Equal(t, *in.Bcc, *got.Bcc)
	assert.Equal(t, in.Subject, got.Subject)
	assert.Equal(t, *in.Text, *got.Text)
	assert.Equal(t, *in.Html, *got.Html)
	assert.Equal(t, *in.ListUnsubscribePost, *got.ListUnsubscribePost)
	assert.Equal(t, *in.ListUnsubscribeUrl, *got.ListUnsubscribeUrl)
	assert.True(t, in.CreatedAt.Equal(got.CreatedAt))
}

func TestMailDao_StoreNullableColumns(t *testing.T) {
	d := dao.NewMailDao(testutil.NewTestDB(t))

	require.NoError(t, d.Store(domain.Mail{MessageID: "id", From: "f", Subject: "s", CreatedAt: time.Now()}))

	got, err := d.FindByMessageID("id")
	require.NoError(t, err)
	assert.Nil(t, got.To)
	assert.Nil(t, got.Cc)
	assert.Nil(t, got.Bcc)
	assert.Nil(t, got.Text)
	assert.Nil(t, got.Html)
	assert.Nil(t, got.ListUnsubscribePost)
	assert.Nil(t, got.ListUnsubscribeUrl)
}

func TestMailDao_StoreDuplicateMessageID(t *testing.T) {
	d := dao.NewMailDao(testutil.NewTestDB(t))

	m := newMail("dup", "a@example.com", time.Now())
	require.NoError(t, d.Store(m))
	assert.Error(t, d.Store(m))
}

func TestMailDao_FindByMessageID_NotFound(t *testing.T) {
	d := dao.NewMailDao(testutil.NewTestDB(t))

	_, err := d.FindByMessageID("missing")
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestMailDao_FindAll(t *testing.T) {
	d := dao.NewMailDao(testutil.NewTestDB(t))
	seed(t, d)

	got, err := d.FindAll(10, 0)
	require.NoError(t, err)
	assert.Equal(t, []string{"m3", "m2", "m1"}, ids(got))

	got, err = d.FindAll(2, 1)
	require.NoError(t, err)
	assert.Equal(t, []string{"m2", "m1"}, ids(got))
}

func TestMailDao_FindByTo(t *testing.T) {
	d := dao.NewMailDao(testutil.NewTestDB(t))
	seed(t, d)

	tests := []struct {
		name   string
		to     *string
		limit  int
		offset int
		want   []string
	}{
		{"nil は全件", nil, 10, 0, []string{"m3", "m2", "m1"}},
		{"空文字は全件", ptr(""), 10, 0, []string{"m3", "m2", "m1"}},
		{"全件 + ページング", nil, 1, 1, []string{"m2"}},
		{"宛先で絞り込み", ptr("a@example.com"), 10, 0, []string{"m3", "m1"}},
		{"宛先で絞り込み + ページング", ptr("a@example.com"), 1, 1, []string{"m1"}},
		{"該当なし", ptr("none@example.com"), 10, 0, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.FindByTo(tt.to, tt.limit, tt.offset)
			require.NoError(t, err)
			assert.Equal(t, tt.want, ids(got))
		})
	}
}

func TestMailDao_CountByTo(t *testing.T) {
	d := dao.NewMailDao(testutil.NewTestDB(t))
	seed(t, d)

	tests := []struct {
		name string
		to   *string
		want int
	}{
		{"nil は全件", nil, 3},
		{"空文字は全件", ptr(""), 3},
		{"宛先で絞り込み", ptr("a@example.com"), 2},
		{"該当なし", ptr("none@example.com"), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.CountByTo(tt.to)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMailDao_Delete(t *testing.T) {
	d := dao.NewMailDao(testutil.NewTestDB(t))
	seed(t, d)

	require.NoError(t, d.Delete("m2"))

	got, err := d.FindAll(10, 0)
	require.NoError(t, err)
	assert.Equal(t, []string{"m3", "m1"}, ids(got))

	// 存在しない ID の削除はエラーにならない
	assert.NoError(t, d.Delete("missing"))
}

func TestMailDao_DeleteAll(t *testing.T) {
	d := dao.NewMailDao(testutil.NewTestDB(t))
	seed(t, d)

	require.NoError(t, d.DeleteAll())

	count, err := d.CountByTo(nil)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestMailDao_ErrorsWhenDBClosed(t *testing.T) {
	db := testutil.NewTestDB(t)
	d := dao.NewMailDao(db)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	assert.Error(t, d.Store(newMail("x", "a@example.com", time.Now())))
	_, err = d.FindByMessageID("x")
	assert.Error(t, err)
	_, err = d.FindAll(10, 0)
	assert.Error(t, err)
	_, err = d.FindByTo(nil, 10, 0)
	assert.Error(t, err)
	_, err = d.FindByTo(ptr("a@example.com"), 10, 0)
	assert.Error(t, err)
	_, err = d.CountByTo(nil)
	assert.Error(t, err)
	_, err = d.CountByTo(ptr("a@example.com"))
	assert.Error(t, err)
	assert.Error(t, d.Delete("x"))
	assert.Error(t, d.DeleteAll())
}
