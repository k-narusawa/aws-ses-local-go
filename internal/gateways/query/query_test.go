package query_test

import (
	"aws-ses-local-go/domain"
	"aws-ses-local-go/internal/gateways/dao"
	gq "aws-ses-local-go/internal/gateways/query"
	"aws-ses-local-go/internal/testutil"
	"aws-ses-local-go/usecase/query"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ptr(s string) *string { return &s }

var base = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

func fullMail() domain.Mail {
	return domain.Mail{
		MessageID:           "m1",
		From:                "from@example.com",
		To:                  ptr("a@example.com"),
		Cc:                  ptr("cc@example.com"),
		Bcc:                 ptr("bcc@example.com"),
		Subject:             "subject",
		Text:                ptr("text"),
		Html:                ptr("<p>html</p>"),
		ListUnsubscribePost: ptr("List-Unsubscribe=One-Click"),
		ListUnsubscribeUrl:  ptr("https://example.com/u"),
		CreatedAt:           base,
	}
}

func setup(t *testing.T) *dao.MailDao {
	t.Helper()
	d := dao.NewMailDao(testutil.NewTestDB(t))
	require.NoError(t, d.Store(fullMail()))
	require.NoError(t, d.Store(domain.Mail{MessageID: "m2", From: "f", To: ptr("b@example.com"), Subject: "s2", CreatedAt: base.Add(time.Minute)}))
	return d
}

func assertFullDto(t *testing.T, got query.MailDto) {
	t.Helper()
	assert.Equal(t, "m1", got.MessageID)
	assert.Equal(t, "from@example.com", got.From)
	assert.Equal(t, "a@example.com", *got.Destination.To)
	assert.Equal(t, "cc@example.com", *got.Destination.Cc)
	assert.Equal(t, "bcc@example.com", *got.Destination.Bcc)
	assert.Equal(t, "subject", got.Subject)
	assert.Equal(t, "text", *got.Body.Text)
	assert.Equal(t, "<p>html</p>", *got.Body.Html)
	assert.Equal(t, "List-Unsubscribe=One-Click", *got.ListUnsubscribePost)
	assert.Equal(t, "https://example.com/u", *got.ListUnsubscribeUrl)
	assert.True(t, base.Equal(got.CreatedAt))
}

func TestMailDtoQueryService_FindByMessageID(t *testing.T) {
	s := gq.NewMailDtoQueryService(*setup(t))

	got, err := s.FindByMessageID("m1")
	require.NoError(t, err)
	assertFullDto(t, got)

	_, err = s.FindByMessageID("missing")
	assert.Error(t, err)
}

func TestMailDtoQueryService_FindAll(t *testing.T) {
	s := gq.NewMailDtoQueryService(*setup(t))

	got, err := s.FindAll(10, 0)
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, "m2", got[0].MessageID)
	assertFullDto(t, got[1])
}

func TestMailDtoQueryService_FindByTo(t *testing.T) {
	s := gq.NewMailDtoQueryService(*setup(t))

	got, err := s.FindByTo(ptr("a@example.com"), 10, 0)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assertFullDto(t, got[0])

	got, err = s.FindByTo(ptr("none@example.com"), 10, 0)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestMailDtoQueryService_Errors(t *testing.T) {
	db := testutil.NewTestDB(t)
	s := gq.NewMailDtoQueryService(*dao.NewMailDao(db))
	sqlDB, _ := db.DB()
	require.NoError(t, sqlDB.Close())

	_, err := s.FindAll(10, 0)
	assert.Error(t, err)
	_, err = s.FindByTo(nil, 10, 0)
	assert.Error(t, err)
}

func TestCountQueryService_CountByTo(t *testing.T) {
	s := gq.NewCountQueryService(*setup(t))

	n, err := s.CountByTo(nil)
	require.NoError(t, err)
	assert.Equal(t, 2, n)

	n, err = s.CountByTo(ptr("a@example.com"))
	require.NoError(t, err)
	assert.Equal(t, 1, n)
}

func TestCountQueryService_Error(t *testing.T) {
	db := testutil.NewTestDB(t)
	s := gq.NewCountQueryService(*dao.NewMailDao(db))
	sqlDB, _ := db.DB()
	require.NoError(t, sqlDB.Close())

	_, err := s.CountByTo(nil)
	assert.Error(t, err)
}
