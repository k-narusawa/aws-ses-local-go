package v1_test

import (
	"aws-ses-local-go/internal/testutil"
	v1 "aws-ses-local-go/usecase/aws/v1"
	"encoding/base64"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_SendEmail(t *testing.T) {
	repo := &testutil.FakeMailRepo{}
	s := v1.NewService(repo)

	out, err := s.SendEmail(v1.SendEmailInput{
		Source:       "from@example.com",
		ToAddresses:  "to@example.com",
		CcAddresses:  "cc@example.com",
		BccAddresses: "bcc@example.com",
		SubjectData:  "subject",
		TextData:     "text",
		HtmlData:     "<p>html</p>",
	})
	require.NoError(t, err)
	require.Len(t, repo.Stored, 1)

	m := repo.Stored[0]
	assert.Equal(t, m.MessageID, out.MessageID)
	assert.Equal(t, "from@example.com", m.From)
	assert.Equal(t, "to@example.com", *m.To)
	assert.Equal(t, "cc@example.com", *m.Cc)
	assert.Equal(t, "bcc@example.com", *m.Bcc)
	assert.Equal(t, "subject", m.Subject)
	assert.Equal(t, "text", *m.Text)
	assert.Equal(t, "<p>html</p>", *m.Html)
	assert.Nil(t, m.ListUnsubscribePost)
	assert.Nil(t, m.ListUnsubscribeUrl)
}

func TestService_SendEmail_StoreError(t *testing.T) {
	s := v1.NewService(&testutil.FakeMailRepo{Err: errors.New("db down")})

	out, err := s.SendEmail(v1.SendEmailInput{Source: "f"})
	assert.EqualError(t, err, "db down")
	assert.Nil(t, out)
}

func TestService_SendRawEmail(t *testing.T) {
	repo := &testutil.FakeMailRepo{}
	s := v1.NewService(repo)

	raw := base64.StdEncoding.EncodeToString([]byte("From: from@example.com\r\nTo: to@example.com\r\nSubject: raw\r\nContent-Type: text/plain\r\n\r\nraw body"))
	out, err := s.SendRawEmail(v1.SendRawEmailInput{RawMessage: raw, Source: "ignored@example.com"})
	require.NoError(t, err)
	require.Len(t, repo.Stored, 1)

	m := repo.Stored[0]
	assert.Equal(t, m.MessageID, out.MessageID)
	// v1 の SendRawEmail は Source ではなくヘッダの From を使う
	assert.Equal(t, "from@example.com", m.From)
	assert.Equal(t, "to@example.com", *m.To)
	assert.Equal(t, "raw", m.Subject)
	assert.Equal(t, "raw body", *m.Text)
}

func TestService_SendRawEmail_Errors(t *testing.T) {
	t.Run("不正な raw message", func(t *testing.T) {
		repo := &testutil.FakeMailRepo{}
		out, err := v1.NewService(repo).SendRawEmail(v1.SendRawEmailInput{RawMessage: "not base64!"})
		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Empty(t, repo.Stored)
	})

	t.Run("保存失敗", func(t *testing.T) {
		raw := base64.StdEncoding.EncodeToString([]byte("From: f\r\nContent-Type: text/plain\r\n\r\nb"))
		out, err := v1.NewService(&testutil.FakeMailRepo{Err: errors.New("db down")}).SendRawEmail(v1.SendRawEmailInput{RawMessage: raw})
		assert.EqualError(t, err, "db down")
		assert.Nil(t, out)
	})
}
