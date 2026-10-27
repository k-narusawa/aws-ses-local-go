package v2_test

import (
	"aws-ses-local-go/internal/testutil"
	v2 "aws-ses-local-go/usecase/aws/v2"
	"encoding/base64"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ptr(s string) *string { return &s }

func simpleInput(body v2.Body) v2.V2EmailOutboundEmailInput {
	return v2.V2EmailOutboundEmailInput{
		FromEmailAddress: ptr("from@example.com"),
		Destination: v2.Destination{
			ToAddresses:  []string{"to1@example.com", "to2@example.com"},
			CcAddresses:  []string{"cc@example.com"},
			BccAddresses: nil,
		},
		Content: v2.Content{
			Simple: &v2.Simple{
				Subject: v2.Subject{Data: "subject"},
				Body:    body,
			},
		},
	}
}

func TestService_SendSimpleEmail(t *testing.T) {
	tests := []struct {
		name     string
		body     v2.Body
		wantText *string
		wantHtml *string
	}{
		{"text のみ", v2.Body{Text: &v2.Text{Data: "text"}}, ptr("text"), nil},
		{"html のみ", v2.Body{Html: &v2.Html{Data: "<p>html</p>"}}, nil, ptr("<p>html</p>")},
		// 現状は text があると html は保存されない
		{"text と html 両方", v2.Body{Text: &v2.Text{Data: "text"}, Html: &v2.Html{Data: "<p>html</p>"}}, ptr("text"), nil},
		{"本文なし", v2.Body{}, nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &testutil.FakeMailRepo{}
			out, err := v2.NewService(repo).SendSimpleEmail(simpleInput(tt.body))
			require.NoError(t, err)
			require.Len(t, repo.Stored, 1)

			m := repo.Stored[0]
			assert.Equal(t, m.MessageID, out.MessageID)
			assert.Equal(t, "from@example.com", m.From)
			assert.Equal(t, "to1@example.com,to2@example.com", *m.To)
			assert.Equal(t, "cc@example.com", *m.Cc)
			assert.Equal(t, "", *m.Bcc)
			assert.Equal(t, "subject", m.Subject)
			assert.Equal(t, tt.wantText, m.Text)
			assert.Equal(t, tt.wantHtml, m.Html)
		})
	}
}

func TestService_SendSimpleEmail_StoreError(t *testing.T) {
	out, err := v2.NewService(&testutil.FakeMailRepo{Err: errors.New("db down")}).
		SendSimpleEmail(simpleInput(v2.Body{Text: &v2.Text{Data: "t"}}))
	assert.EqualError(t, err, "db down")
	assert.Nil(t, out)
}

func rawInput(from *string, data string) v2.V2EmailOutboundEmailInput {
	return v2.V2EmailOutboundEmailInput{
		FromEmailAddress: from,
		Content:          v2.Content{Raw: &v2.Raw{Data: data}},
	}
}

func TestService_SendRawEmail(t *testing.T) {
	raw := base64.StdEncoding.EncodeToString([]byte("From: header@example.com\r\nTo: to@example.com\r\nSubject: raw\r\nContent-Type: text/plain\r\n\r\nraw body"))

	t.Run("FromEmailAddress を優先する", func(t *testing.T) {
		repo := &testutil.FakeMailRepo{}
		out, err := v2.NewService(repo).SendRawEmail(rawInput(ptr("from@example.com"), raw))
		require.NoError(t, err)
		require.Len(t, repo.Stored, 1)
		assert.Equal(t, repo.Stored[0].MessageID, out.MessageID)
		assert.Equal(t, "from@example.com", repo.Stored[0].From)
		assert.Equal(t, "raw body", *repo.Stored[0].Text)
	})

	t.Run("FromEmailAddress が nil ならヘッダの From", func(t *testing.T) {
		repo := &testutil.FakeMailRepo{}
		_, err := v2.NewService(repo).SendRawEmail(rawInput(nil, raw))
		require.NoError(t, err)
		assert.Equal(t, "header@example.com", repo.Stored[0].From)
	})
}

func TestService_SendRawEmail_Errors(t *testing.T) {
	t.Run("不正な raw message", func(t *testing.T) {
		repo := &testutil.FakeMailRepo{}
		out, err := v2.NewService(repo).SendRawEmail(rawInput(nil, "not base64!"))
		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Empty(t, repo.Stored)
	})

	t.Run("保存失敗", func(t *testing.T) {
		raw := base64.StdEncoding.EncodeToString([]byte("From: f\r\nContent-Type: text/plain\r\n\r\nb"))
		out, err := v2.NewService(&testutil.FakeMailRepo{Err: errors.New("db down")}).SendRawEmail(rawInput(nil, raw))
		assert.EqualError(t, err, "db down")
		assert.Nil(t, out)
	})
}
