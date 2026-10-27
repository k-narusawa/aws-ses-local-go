// Package e2e は実際の AWS SDK (v1 / v2) クライアントからローカルサーバーへ送信し、
// SQLite への保存と /emails API での取得までを通しで検証する。
package e2e_test

import (
	"aws-ses-local-go/internal/controllers/rest"
	"aws-ses-local-go/internal/gateways/dao"
	"aws-ses-local-go/internal/gateways/query"
	"aws-ses-local-go/internal/gateways/repository"
	"aws-ses-local-go/internal/testutil"
	v1 "aws-ses-local-go/usecase/aws/v1"
	v2 "aws-ses-local-go/usecase/aws/v2"
	"aws-ses-local-go/usecase/mail"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	sesv2 "github.com/aws/aws-sdk-go-v2/service/sesv2"
	sesv2types "github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// app/main.go と同じ組み立てでサーバーを起動する
func startServer(t *testing.T) *httptest.Server {
	t.Helper()

	db := testutil.NewTestDB(t)
	mailDao := dao.NewMailDao(db)
	mailRepo := repository.NewMailRepository(*mailDao)

	e := echo.New()
	rest.NewAwsHandler(e, v1.NewService(mailRepo), v2.NewService(mailRepo))
	rest.NewMailHandler(e,
		query.NewMailDtoQueryService(*mailDao),
		query.NewCountQueryService(*mailDao),
		mail.NewService(mailRepo),
	)

	srv := httptest.NewServer(e)
	t.Cleanup(srv.Close)
	return srv
}

func newV1Client(t *testing.T, endpoint string) *ses.SES {
	t.Helper()
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Endpoint:    aws.String(endpoint),
		Credentials: credentials.NewStaticCredentials("dummy", "dummy", ""),
	})
	require.NoError(t, err)
	return ses.New(sess)
}

func newV2Client(endpoint string) *sesv2.Client {
	return sesv2.New(sesv2.Options{
		Region:       "us-west-2",
		BaseEndpoint: aws.String(endpoint),
	})
}

type mailsResponse struct {
	TotalSize int `json:"total_size"`
	Items     []struct {
		MessageID   string `json:"message_id"`
		From        string `json:"from"`
		Subject     string `json:"subject"`
		Destination struct {
			To  *string `json:"to"`
			Cc  *string `json:"cc"`
			Bcc *string `json:"bcc"`
		} `json:"destination"`
		Body struct {
			Text *string `json:"text"`
			Html *string `json:"html"`
		} `json:"body"`
		ListUnsubscribePost *string `json:"list_unsubscribe_post"`
		ListUnsubscribeUrl  *string `json:"list_unsubscribe_url"`
	} `json:"items"`
}

func getMails(t *testing.T, base string) mailsResponse {
	t.Helper()
	res, err := http.Get(base + "/emails")
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, 200, res.StatusCode)

	var out mailsResponse
	require.NoError(t, json.NewDecoder(res.Body).Decode(&out))
	return out
}

func skipIfShort(t *testing.T) {
	if testing.Short() {
		t.Skip("e2e test is skipped in short mode")
	}
}

func TestE2E_V1SendEmail(t *testing.T) {
	skipIfShort(t)
	srv := startServer(t)
	client := newV1Client(t, srv.URL)

	out, err := client.SendEmail(&ses.SendEmailInput{
		Source: aws.String("from@example.com"),
		Destination: &ses.Destination{
			ToAddresses:  []*string{aws.String("to@example.com")},
			CcAddresses:  []*string{aws.String("cc@example.com")},
			BccAddresses: []*string{aws.String("bcc@example.com")},
		},
		Message: &ses.Message{
			Subject: &ses.Content{Data: aws.String("v1 件名"), Charset: aws.String("UTF-8")},
			Body: &ses.Body{
				Text: &ses.Content{Data: aws.String("v1 本文"), Charset: aws.String("UTF-8")},
				Html: &ses.Content{Data: aws.String("<p>v1</p>"), Charset: aws.String("UTF-8")},
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, out.MessageId)

	got := getMails(t, srv.URL)
	require.Equal(t, 1, got.TotalSize)
	m := got.Items[0]
	assert.Equal(t, *out.MessageId, m.MessageID)
	assert.Equal(t, "from@example.com", m.From)
	assert.Equal(t, "v1 件名", m.Subject)
	assert.Equal(t, "to@example.com", *m.Destination.To)
	assert.Equal(t, "cc@example.com", *m.Destination.Cc)
	assert.Equal(t, "bcc@example.com", *m.Destination.Bcc)
	assert.Equal(t, "v1 本文", *m.Body.Text)
	assert.Equal(t, "<p>v1</p>", *m.Body.Html)
}

func TestE2E_V1SendRawEmail(t *testing.T) {
	skipIfShort(t)
	srv := startServer(t)
	client := newV1Client(t, srv.URL)

	raw := strings.Join([]string{
		"From: raw-from@example.com",
		"To: raw-to@example.com",
		"Subject: v1 raw",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		"v1 raw body",
	}, "\r\n")

	out, err := client.SendRawEmail(&ses.SendRawEmailInput{
		RawMessage: &ses.RawMessage{Data: []byte(raw)},
	})
	require.NoError(t, err)

	got := getMails(t, srv.URL)
	require.Equal(t, 1, got.TotalSize)
	m := got.Items[0]
	assert.Equal(t, *out.MessageId, m.MessageID)
	assert.Equal(t, "raw-from@example.com", m.From)
	assert.Equal(t, "raw-to@example.com", *m.Destination.To)
	assert.Equal(t, "v1 raw", m.Subject)
	assert.Equal(t, "v1 raw body", *m.Body.Text)
}

func TestE2E_V2SendSimpleEmail(t *testing.T) {
	skipIfShort(t)
	srv := startServer(t)
	client := newV2Client(srv.URL)

	out, err := client.SendEmail(context.Background(), &sesv2.SendEmailInput{
		FromEmailAddress: aws.String("from@example.com"),
		Destination: &sesv2types.Destination{
			ToAddresses: []string{"to1@example.com", "to2@example.com"},
		},
		Content: &sesv2types.EmailContent{
			Simple: &sesv2types.Message{
				Subject: &sesv2types.Content{Data: aws.String("v2 件名")},
				Body:    &sesv2types.Body{Text: &sesv2types.Content{Data: aws.String("v2 本文")}},
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, out.MessageId)

	got := getMails(t, srv.URL)
	require.Equal(t, 1, got.TotalSize)
	m := got.Items[0]
	assert.Equal(t, *out.MessageId, m.MessageID)
	assert.Equal(t, "from@example.com", m.From)
	assert.Equal(t, "to1@example.com,to2@example.com", *m.Destination.To)
	assert.Equal(t, "v2 件名", m.Subject)
	assert.Equal(t, "v2 本文", *m.Body.Text)
}

func TestE2E_V2SendRawEmail(t *testing.T) {
	skipIfShort(t)
	srv := startServer(t)
	client := newV2Client(srv.URL)

	raw := strings.Join([]string{
		"From: header@example.com",
		"To: to@example.com",
		"Subject: =?UTF-8?B?5pel5pys6Kqe44Gu44Oh44O844Or?=",
		"List-Unsubscribe: <https://example.com/unsubscribe>",
		"List-Unsubscribe-Post: List-Unsubscribe=One-Click",
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"Content-Transfer-Encoding: base64",
		"",
		"5pel5pys6Kqe44Gu5pys5paH",
	}, "\r\n")

	out, err := client.SendEmail(context.Background(), &sesv2.SendEmailInput{
		FromEmailAddress: aws.String("from@example.com"),
		Content: &sesv2types.EmailContent{
			Raw: &sesv2types.RawMessage{Data: []byte(raw)},
		},
	})
	require.NoError(t, err)

	got := getMails(t, srv.URL)
	require.Equal(t, 1, got.TotalSize)
	m := got.Items[0]
	assert.Equal(t, *out.MessageId, m.MessageID)
	assert.Equal(t, "from@example.com", m.From)
	assert.Equal(t, "to@example.com", *m.Destination.To)
	assert.Equal(t, "日本語のメール", m.Subject)
	assert.Equal(t, "日本語の本文", *m.Body.Text)
	assert.Equal(t, "https://example.com/unsubscribe", *m.ListUnsubscribeUrl)
	assert.Equal(t, "List-Unsubscribe=One-Click", *m.ListUnsubscribePost)
}

func TestE2E_DeleteAndFilter(t *testing.T) {
	skipIfShort(t)
	srv := startServer(t)
	client := newV2Client(srv.URL)

	send := func(to string) string {
		out, err := client.SendEmail(context.Background(), &sesv2.SendEmailInput{
			FromEmailAddress: aws.String("from@example.com"),
			Destination:      &sesv2types.Destination{ToAddresses: []string{to}},
			Content: &sesv2types.EmailContent{Simple: &sesv2types.Message{
				Subject: &sesv2types.Content{Data: aws.String("s")},
				Body:    &sesv2types.Body{Text: &sesv2types.Content{Data: aws.String("b")}},
			}},
		})
		require.NoError(t, err)
		return *out.MessageId
	}
	id1 := send("a@example.com")
	send("a@example.com")
	send("b@example.com")

	res, err := http.Get(srv.URL + "/emails?to_address=a%40example.com")
	require.NoError(t, err)
	var filtered mailsResponse
	require.NoError(t, json.NewDecoder(res.Body).Decode(&filtered))
	res.Body.Close()
	assert.Equal(t, 2, filtered.TotalSize)

	doDelete := func(path string) int {
		req, err := http.NewRequest(http.MethodDelete, srv.URL+path, nil)
		require.NoError(t, err)
		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		res.Body.Close()
		return res.StatusCode
	}

	assert.Equal(t, 204, doDelete(fmt.Sprintf("/emails/%s", id1)))
	assert.Equal(t, 2, getMails(t, srv.URL).TotalSize)

	assert.Equal(t, 204, doDelete("/emails"))
	assert.Equal(t, 0, getMails(t, srv.URL).TotalSize)
}
