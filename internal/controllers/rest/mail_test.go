package rest_test

import (
	"aws-ses-local-go/internal/controllers/rest"
	"aws-ses-local-go/usecase/query"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeMailQuery struct {
	mails     []query.MailDto
	err       error
	gotTo     *string
	gotLimit  int
	gotOffset int
}

func (f *fakeMailQuery) FindAll(limit, offset int) ([]query.MailDto, error) {
	return f.mails, f.err
}

func (f *fakeMailQuery) FindByMessageID(string) (query.MailDto, error) {
	return query.MailDto{}, f.err
}

func (f *fakeMailQuery) FindByTo(to *string, limit, offset int) ([]query.MailDto, error) {
	f.gotTo, f.gotLimit, f.gotOffset = to, limit, offset
	return f.mails, f.err
}

type fakeCountQuery struct {
	count int
	err   error
	gotTo *string
}

func (f *fakeCountQuery) CountByTo(to *string) (int, error) {
	f.gotTo = to
	return f.count, f.err
}

type fakeMailService struct {
	deletedID  string
	deletedAll bool
	err        error
}

func (f *fakeMailService) DeleteMail(id string) error {
	f.deletedID = id
	return f.err
}

func (f *fakeMailService) DeleteMails() error {
	f.deletedAll = true
	return f.err
}

func newMailServer(q *fakeMailQuery, c *fakeCountQuery, s *fakeMailService) *echo.Echo {
	e := echo.New()
	rest.NewMailHandler(e, q, c, s)
	return e
}

func doRequest(e *echo.Echo, method, target string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(method, target, nil))
	return rec
}

func TestMailHandler_GetMails(t *testing.T) {
	to := "to@example.com"
	created := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)
	mails := []query.MailDto{{
		MessageID:   "m1",
		From:        "from@example.com",
		Destination: query.Destination{To: &to},
		Subject:     "subject",
		Body:        query.Body{Text: &to},
		CreatedAt:   created,
	}}

	tests := []struct {
		name          string
		target        string
		count         int
		wantTo        string
		wantLimit     int
		wantOffset    int
		wantPage      int
		wantTotalPage int
	}{
		{"デフォルト", "/emails", 25, "", 10, 0, 1, 2},
		{"ページ・件数指定", "/emails?page=3&limit=5", 25, "", 5, 10, 3, 5},
		{"宛先指定", "/emails?to_address=to%40example.com", 1, "to@example.com", 10, 0, 1, 0},
		{"/store も同じハンドラ", "/store?page=2&limit=1", 3, "", 1, 1, 2, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &fakeMailQuery{mails: mails}
			c := &fakeCountQuery{count: tt.count}
			rec := doRequest(newMailServer(q, c, &fakeMailService{}), http.MethodGet, tt.target)

			require.Equal(t, 200, rec.Code)
			assert.Equal(t, tt.wantTo, *q.gotTo)
			assert.Equal(t, tt.wantTo, *c.gotTo)
			assert.Equal(t, tt.wantLimit, q.gotLimit)
			assert.Equal(t, tt.wantOffset, q.gotOffset)

			var resp rest.MailResponse
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			assert.Equal(t, tt.wantPage, resp.Page)
			assert.Equal(t, tt.wantLimit, resp.Limit)
			assert.Equal(t, 1, resp.Size)
			assert.Equal(t, tt.wantTotalPage, resp.TotalPage)
			assert.Equal(t, tt.count, resp.TotalSize)
			require.Len(t, resp.Items, 1)
			assert.Equal(t, "m1", resp.Items[0].MessageID)
			assert.True(t, created.Equal(resp.Items[0].CreatedAt))
		})
	}
}

func TestMailHandler_GetMails_JSONShape(t *testing.T) {
	// フロントエンドが依存している JSON のキー名を固定する
	to := "to@example.com"
	q := &fakeMailQuery{mails: []query.MailDto{{
		MessageID:   "m1",
		From:        "from@example.com",
		Destination: query.Destination{To: &to},
		Subject:     "subject",
		CreatedAt:   time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC),
	}}}
	rec := doRequest(newMailServer(q, &fakeCountQuery{count: 1}, &fakeMailService{}), http.MethodGet, "/emails")

	require.Equal(t, 200, rec.Code)
	assert.JSONEq(t, `{
		"page": 1, "limit": 10, "size": 1, "total_page": 0, "total_size": 1,
		"items": [{
			"message_id": "m1",
			"from": "from@example.com",
			"destination": {"to": "to@example.com", "cc": null, "bcc": null},
			"subject": "subject",
			"body": {"text": null, "html": null},
			"list_unsubscribe_post": null,
			"list_unsubscribe_url": null,
			"created_at": "2025-01-02T03:04:05Z"
		}]
	}`, rec.Body.String())
}

func TestMailHandler_GetMails_EmptyIsArray(t *testing.T) {
	rec := doRequest(newMailServer(&fakeMailQuery{mails: nil}, &fakeCountQuery{}, &fakeMailService{}), http.MethodGet, "/emails")

	require.Equal(t, 200, rec.Code)
	var raw map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &raw))
	assert.Equal(t, "[]", string(raw["items"]))
}

func TestMailHandler_GetMails_Errors(t *testing.T) {
	t.Run("メール取得失敗", func(t *testing.T) {
		rec := doRequest(newMailServer(&fakeMailQuery{err: errors.New("boom")}, &fakeCountQuery{}, &fakeMailService{}), http.MethodGet, "/emails")
		assert.Equal(t, 500, rec.Code)
	})

	t.Run("件数取得失敗", func(t *testing.T) {
		rec := doRequest(newMailServer(&fakeMailQuery{}, &fakeCountQuery{err: errors.New("boom")}, &fakeMailService{}), http.MethodGet, "/emails")
		assert.Equal(t, 500, rec.Code)
	})
}

func TestMailHandler_DeleteMail(t *testing.T) {
	s := &fakeMailService{}
	rec := doRequest(newMailServer(&fakeMailQuery{}, &fakeCountQuery{}, s), http.MethodDelete, "/emails/ses-123")
	assert.Equal(t, 204, rec.Code)
	assert.Equal(t, "ses-123", s.deletedID)

	rec = doRequest(newMailServer(&fakeMailQuery{}, &fakeCountQuery{}, &fakeMailService{err: errors.New("boom")}), http.MethodDelete, "/emails/ses-123")
	assert.Equal(t, 500, rec.Code)
}

func TestMailHandler_DeleteMails(t *testing.T) {
	s := &fakeMailService{}
	rec := doRequest(newMailServer(&fakeMailQuery{}, &fakeCountQuery{}, s), http.MethodDelete, "/emails")
	assert.Equal(t, 204, rec.Code)
	assert.True(t, s.deletedAll)

	rec = doRequest(newMailServer(&fakeMailQuery{}, &fakeCountQuery{}, &fakeMailService{err: errors.New("boom")}), http.MethodDelete, "/emails")
	assert.Equal(t, 500, rec.Code)
}
