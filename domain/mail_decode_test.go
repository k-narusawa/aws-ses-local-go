package domain_test

import (
	"aws-ses-local-go/domain"
	"encoding/base64"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/encoding/japanese"
)

var messageIDPattern = regexp.MustCompile(`^ses-[1-9][0-9]{8}$`)

func encodeRaw(lines ...string) string {
	return base64.StdEncoding.EncodeToString([]byte(strings.Join(lines, "\r\n")))
}

func ptr(s string) *string { return &s }

func TestNewMail(t *testing.T) {
	before := time.Now()
	m := domain.NewMail(
		"from@example.com",
		ptr("to@example.com"),
		ptr("cc@example.com"),
		ptr("bcc@example.com"),
		"subject",
		ptr("text body"),
		ptr("<p>html body</p>"),
		ptr("List-Unsubscribe=One-Click"),
		ptr("https://example.com/unsubscribe"),
	)
	after := time.Now()

	assert.Regexp(t, messageIDPattern, m.MessageID)
	assert.Equal(t, "from@example.com", m.From)
	assert.Equal(t, "to@example.com", *m.To)
	assert.Equal(t, "cc@example.com", *m.Cc)
	assert.Equal(t, "bcc@example.com", *m.Bcc)
	assert.Equal(t, "subject", m.Subject)
	assert.Equal(t, "text body", *m.Text)
	assert.Equal(t, "<p>html body</p>", *m.Html)
	assert.Equal(t, "List-Unsubscribe=One-Click", *m.ListUnsubscribePost)
	assert.Equal(t, "https://example.com/unsubscribe", *m.ListUnsubscribeUrl)
	assert.False(t, m.CreatedAt.Before(before))
	assert.False(t, m.CreatedAt.After(after))
}

func TestNewMail_MessageIDIsUnique(t *testing.T) {
	seen := map[string]struct{}{}
	for i := 0; i < 100; i++ {
		m := domain.NewMail("f", nil, nil, nil, "s", nil, nil, nil, nil)
		assert.Regexp(t, messageIDPattern, m.MessageID)
		seen[m.MessageID] = struct{}{}
	}
	// 9 桁の乱数なので 100 件程度で衝突することはまずない
	assert.Greater(t, len(seen), 95)
}

func TestMail_TableName(t *testing.T) {
	assert.Equal(t, "mails", domain.Mail{}.TableName())
}

func TestFromRawEmailRequest_Decoding(t *testing.T) {
	tests := []struct {
		name        string
		raw         string
		from        *string
		wantFrom    string
		wantTo      string
		wantSubject string
		wantText    string
	}{
		{
			name: "7bit で Content-Transfer-Encoding 省略",
			raw: encodeRaw(
				"From: sender@example.com",
				"To: recipient@example.com",
				"Subject: Plain",
				"Content-Type: text/plain; charset=utf-8",
				"",
				"plain body",
			),
			wantFrom:    "sender@example.com",
			wantTo:      "recipient@example.com",
			wantSubject: "Plain",
			wantText:    "plain body",
		},
		{
			name: "charset 指定なしは utf-8 として扱う",
			raw: encodeRaw(
				"From: sender@example.com",
				"To: recipient@example.com",
				"Subject: NoCharset",
				"Content-Type: text/plain",
				"",
				"こんにちは",
			),
			wantFrom:    "sender@example.com",
			wantTo:      "recipient@example.com",
			wantSubject: "NoCharset",
			wantText:    "こんにちは",
		},
		{
			name: "base64 本文と MIME エンコードされた件名 (UTF-8)",
			raw: encodeRaw(
				"From: sender@example.com",
				"To: recipient@example.com",
				"Subject: =?UTF-8?B?"+base64.StdEncoding.EncodeToString([]byte("日本語の件名"))+"?=",
				"Content-Type: text/plain; charset=UTF-8",
				"Content-Transfer-Encoding: base64",
				"",
				base64.StdEncoding.EncodeToString([]byte("日本語の本文")),
			),
			wantFrom:    "sender@example.com",
			wantTo:      "recipient@example.com",
			wantSubject: "日本語の件名",
			wantText:    "日本語の本文",
		},
		{
			name: "quoted-printable 本文と Q エンコードされた件名",
			raw: encodeRaw(
				"From: sender@example.com",
				"To: recipient@example.com",
				"Subject: =?UTF-8?Q?caf=C3=A9?=",
				"Content-Type: text/plain; charset=utf-8",
				"Content-Transfer-Encoding: Quoted-Printable",
				"",
				"caf=C3=A9 =",
				"au lait",
			),
			wantFrom:    "sender@example.com",
			wantTo:      "recipient@example.com",
			wantSubject: "café",
			wantText:    "café au lait",
		},
		{
			name: "FromEmailAddress が指定されたらヘッダより優先する",
			raw: encodeRaw(
				"From: header@example.com",
				"To: recipient@example.com",
				"Subject: Override",
				"Content-Type: text/plain",
				"",
				"body",
			),
			from:        ptr("override@example.com"),
			wantFrom:    "override@example.com",
			wantTo:      "recipient@example.com",
			wantSubject: "Override",
			wantText:    "body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := domain.FromRawEmailRequest(tt.from, tt.raw)
			require.NoError(t, err)

			assert.Regexp(t, messageIDPattern, m.MessageID)
			assert.Equal(t, tt.wantFrom, m.From)
			require.NotNil(t, m.To)
			assert.Equal(t, tt.wantTo, *m.To)
			assert.Equal(t, tt.wantSubject, m.Subject)
			require.NotNil(t, m.Text)
			assert.Equal(t, tt.wantText, *m.Text)
			assert.False(t, m.CreatedAt.IsZero())
		})
	}
}

func TestFromRawEmailRequest_ListUnsubscribe(t *testing.T) {
	raw := encodeRaw(
		"From: sender@example.com",
		"To: recipient@example.com",
		"Subject: Unsubscribe",
		"List-Unsubscribe: <https://example.com/unsubscribe?id=1>",
		"List-Unsubscribe-Post: List-Unsubscribe=One-Click",
		"Content-Type: text/plain",
		"",
		"body",
	)

	m, err := domain.FromRawEmailRequest(nil, raw)
	require.NoError(t, err)
	require.NotNil(t, m.ListUnsubscribeUrl)
	require.NotNil(t, m.ListUnsubscribePost)
	assert.Equal(t, "https://example.com/unsubscribe?id=1", *m.ListUnsubscribeUrl)
	assert.Equal(t, "List-Unsubscribe=One-Click", *m.ListUnsubscribePost)
}

func TestFromRawEmailRequest_Errors(t *testing.T) {
	tests := []struct {
		name string
		raw  string
	}{
		{
			name: "Content-Type ヘッダがない",
			raw: encodeRaw(
				"From: sender@example.com",
				"Subject: No content type",
				"",
				"body",
			),
		},
		{
			name: "Content-Type が不正",
			raw: encodeRaw(
				"From: sender@example.com",
				"Content-Type: ;;;",
				"",
				"body",
			),
		},
		{
			name: "本文の base64 が壊れている",
			raw: encodeRaw(
				"From: sender@example.com",
				"Content-Type: text/plain; charset=utf-8",
				"Content-Transfer-Encoding: base64",
				"",
				"!!!not-base64!!!",
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.FromRawEmailRequest(nil, tt.raw)
			assert.Error(t, err)
		})
	}
}

func TestFromRawEmailRequest_InvalidEncodedSubjectIsIgnored(t *testing.T) {
	// 件名のデコードに失敗してもエラーにはせず空文字になる (現在の挙動)
	raw := encodeRaw(
		"From: sender@example.com",
		"To: recipient@example.com",
		"Subject: =?UNKNOWN-CHARSET?B?AAAA?=",
		"Content-Type: text/plain",
		"",
		"body",
	)

	m, err := domain.FromRawEmailRequest(nil, raw)
	require.NoError(t, err)
	assert.Equal(t, "", m.Subject)
	assert.Equal(t, "body", *m.Text)
}

// 現在の実装は charset.NewReader に Content-Type ではなく charset 名を渡しているため、
// charset パラメータが実質無視される。また件名の WordDecoder に CharsetReader を設定していないため
// UTF-8 / ISO-8859-1 / US-ASCII 以外の件名はデコードに失敗して空になる。
// Go やライブラリのアップグレードで挙動が変わったことを検知するため、現状の挙動を固定しておく。
// (バグを修正した場合はこのテストを期待値ごと更新すること)
func TestFromRawEmailRequest_NonUTF8Charset_CurrentBehavior(t *testing.T) {
	body := mustISO2022JP(t, "日本語の本文です")
	raw := encodeRaw(
		"From: sender@example.com",
		"To: recipient@example.com",
		"Subject: =?ISO-2022-JP?B?"+base64.StdEncoding.EncodeToString([]byte(mustISO2022JP(t, "件名")))+"?=",
		"Content-Type: text/plain; charset=ISO-2022-JP",
		"Content-Transfer-Encoding: 7bit",
		"",
		body,
	)

	m, err := domain.FromRawEmailRequest(nil, raw)
	require.NoError(t, err)
	assert.Equal(t, "", m.Subject)
	assert.Equal(t, body, *m.Text)
}

func TestFromRawEmailRequest_UnknownCharset_CurrentBehavior(t *testing.T) {
	// 上記と同じ理由で、未知の charset でもエラーにならず本文はそのまま返る (現状の挙動)
	raw := encodeRaw(
		"From: sender@example.com",
		"Content-Type: text/plain; charset=unknown-charset-xyz",
		"",
		"body",
	)

	m, err := domain.FromRawEmailRequest(nil, raw)
	require.NoError(t, err)
	assert.Equal(t, "body", *m.Text)
}

func mustISO2022JP(t *testing.T, s string) string {
	t.Helper()
	out, err := japanese.ISO2022JP.NewEncoder().String(s)
	require.NoError(t, err)
	return out
}
