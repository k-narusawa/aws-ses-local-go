package rest_test

import (
	"aws-ses-local-go/internal/controllers/rest"
	v1 "aws-ses-local-go/usecase/aws/v1"
	v2 "aws-ses-local-go/usecase/aws/v2"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeV1 struct {
	sendIn    *v1.SendEmailInput
	rawIn     *v1.SendRawEmailInput
	err       error
	messageID string
}

func (f *fakeV1) SendEmail(in v1.SendEmailInput) (*v1.SendEmailOutput, error) {
	f.sendIn = &in
	if f.err != nil {
		return nil, f.err
	}
	return &v1.SendEmailOutput{MessageID: f.messageID}, nil
}

func (f *fakeV1) SendRawEmail(in v1.SendRawEmailInput) (*v1.SendEmailOutput, error) {
	f.rawIn = &in
	if f.err != nil {
		return nil, f.err
	}
	return &v1.SendEmailOutput{MessageID: f.messageID}, nil
}

type fakeV2 struct {
	simpleIn  *v2.V2EmailOutboundEmailInput
	rawIn     *v2.V2EmailOutboundEmailInput
	err       error
	messageID string
}

func (f *fakeV2) SendSimpleEmail(in v2.V2EmailOutboundEmailInput) (*v2.SendEmailV2Output, error) {
	f.simpleIn = &in
	if f.err != nil {
		return nil, f.err
	}
	return &v2.SendEmailV2Output{MessageID: f.messageID}, nil
}

func (f *fakeV2) SendRawEmail(in v2.V2EmailOutboundEmailInput) (*v2.SendEmailV2Output, error) {
	f.rawIn = &in
	if f.err != nil {
		return nil, f.err
	}
	return &v2.SendEmailV2Output{MessageID: f.messageID}, nil
}

func newAwsServer(f1 *fakeV1, f2 *fakeV2) *echo.Echo {
	e := echo.New()
	rest.NewAwsHandler(e, f1, f2)
	return e
}

func postForm(e *echo.Echo, form url.Values) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

func postJSON(e *echo.Echo, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

func TestAwsHandler_V1SendEmail(t *testing.T) {
	f1 := &fakeV1{messageID: "ses-111111111"}
	e := newAwsServer(f1, &fakeV2{})

	rec := postForm(e, url.Values{
		"Action":                            {"SendEmail"},
		"Version":                           {"2010-12-01"},
		"ConfigurationSetName":              {"cfg"},
		"Source":                            {"from@example.com"},
		"Destination.ToAddresses.member.1":  {"to@example.com"},
		"Destination.CcAddresses.member.1":  {"cc@example.com"},
		"Destination.BccAddresses.member.1": {"bcc@example.com"},
		"Message.Subject.Data":              {"subject"},
		"Message.Subject.Charset":           {"UTF-8"},
		"Message.Body.Text.Data":            {"text"},
		"Message.Body.Text.Charset":         {"UTF-8"},
		"Message.Body.Html.Data":            {"<p>html</p>"},
		"Message.Body.Html.Charset":         {"UTF-8"},
		"ReplyToAddresses.member.1":         {"reply@example.com"},
		"ReturnPath":                        {"return@example.com"},
		"ReturnPathArn":                     {"arn:return"},
		"SourceArn":                         {"arn:source"},
		"FromArn":                           {"arn:from"},
		"Tags.member.1":                     {"tag"},
		"Destination.member.1":              {"dest"},
	})

	assert.Equal(t, 200, rec.Code)
	assert.Contains(t, rec.Header().Get(echo.HeaderContentType), "application/xml")
	assert.Contains(t, rec.Body.String(), `<SendEmailResponse xmlns="http://ses.amazonaws.com/doc/2010-12-01/"><SendEmailResult><MessageId>ses-111111111</MessageId></SendEmailResult></SendEmailResponse>`)

	require.NotNil(t, f1.sendIn)
	assert.Equal(t, v1.SendEmailInput{
		Version:              "2010-12-01",
		ConfigurationSetName: "cfg",
		ToAddresses:          "to@example.com",
		CcAddresses:          "cc@example.com",
		BccAddresses:         "bcc@example.com",
		HtmlData:             "<p>html</p>",
		HtmlCharset:          "UTF-8",
		TextData:             "text",
		TextCharset:          "UTF-8",
		SubjectData:          "subject",
		SubjectCharset:       "UTF-8",
		ReplyToAddresses:     "reply@example.com",
		ReturnPath:           "return@example.com",
		ReturnPathArn:        "arn:return",
		Source:               "from@example.com",
		SourceArn:            "arn:source",
		Tags:                 "tag",
		Destination:          "dest",
		FromArn:              "arn:from",
	}, *f1.sendIn)
}

func TestAwsHandler_V1SendRawEmail(t *testing.T) {
	f1 := &fakeV1{messageID: "ses-222222222"}
	e := newAwsServer(f1, &fakeV2{})

	rec := postForm(e, url.Values{
		"Action":          {"SendRawEmail"},
		"Version":         {"2010-12-01"},
		"Source":          {"from@example.com"},
		"RawMessage.Data": {"cmF3"},
	})

	assert.Equal(t, 200, rec.Code)
	assert.Contains(t, rec.Body.String(), `<SendRawEmailResponse xmlns="http://ses.amazonaws.com/doc/2010-12-01/"><SendRawEmailResult><MessageId>ses-222222222</MessageId></SendRawEmailResult></SendRawEmailResponse>`)
	require.NotNil(t, f1.rawIn)
	assert.Equal(t, "cmF3", f1.rawIn.RawMessage)
	assert.Equal(t, "from@example.com", f1.rawIn.Source)
	assert.Equal(t, "2010-12-01", f1.rawIn.Version)
}

func TestAwsHandler_V1Errors(t *testing.T) {
	t.Run("SendEmail 失敗は 500", func(t *testing.T) {
		rec := postForm(newAwsServer(&fakeV1{err: errors.New("boom")}, &fakeV2{}), url.Values{"Action": {"SendEmail"}})
		assert.Equal(t, 500, rec.Code)
	})

	t.Run("SendRawEmail 失敗は 500", func(t *testing.T) {
		rec := postForm(newAwsServer(&fakeV1{err: errors.New("boom")}, &fakeV2{}), url.Values{"Action": {"SendRawEmail"}})
		assert.Equal(t, 500, rec.Code)
	})

	t.Run("未知の Action は 400", func(t *testing.T) {
		f1 := &fakeV1{}
		rec := postForm(newAwsServer(f1, &fakeV2{}), url.Values{"Action": {"GetSendQuota"}})
		assert.Equal(t, 400, rec.Code)
		assert.JSONEq(t, `"Invalid action"`, rec.Body.String())
		assert.Nil(t, f1.sendIn)
		assert.Nil(t, f1.rawIn)
	})
}

func TestAwsHandler_V2Simple(t *testing.T) {
	f2 := &fakeV2{messageID: "ses-333333333"}
	e := newAwsServer(&fakeV1{}, f2)

	rec := postJSON(e, "/v2/email/outbound-emails", `{
		"FromEmailAddress": "from@example.com",
		"Destination": {"ToAddresses": ["to@example.com"], "CcAddresses": ["cc@example.com"]},
		"Content": {"Simple": {
			"Subject": {"Data": "subject", "CharSet": "UTF-8"},
			"Body": {"Text": {"Data": "text"}, "Html": {"Data": "<p>html</p>"}}
		}},
		"EmailTags": [{"Name": "k", "Value": "v"}],
		"ReplyToAddresses": ["reply@example.com"],
		"ListManagementOptions": {"ContactListName": "list", "TopicName": "topic"}
	}`)

	assert.Equal(t, 200, rec.Code)
	assert.JSONEq(t, `{"MessageId":"ses-333333333"}`, rec.Body.String())

	require.NotNil(t, f2.simpleIn)
	in := f2.simpleIn
	assert.Nil(t, f2.rawIn)
	assert.Equal(t, "from@example.com", *in.FromEmailAddress)
	assert.Equal(t, []string{"to@example.com"}, in.Destination.ToAddresses)
	assert.Equal(t, []string{"cc@example.com"}, in.Destination.CcAddresses)
	assert.Equal(t, "subject", in.Content.Simple.Subject.Data)
	assert.Equal(t, "UTF-8", *in.Content.Simple.Subject.CharSet)
	assert.Equal(t, "text", in.Content.Simple.Body.Text.Data)
	assert.Equal(t, "<p>html</p>", in.Content.Simple.Body.Html.Data)
	assert.Equal(t, "k", in.EmailTags[0].Name)
	assert.Equal(t, []string{"reply@example.com"}, in.ReplyToAddresses)
	assert.Equal(t, "topic", *in.ListManagementOptions.TopicName)
}

func TestAwsHandler_V2Raw(t *testing.T) {
	f2 := &fakeV2{messageID: "ses-444444444"}
	e := newAwsServer(&fakeV1{}, f2)

	rec := postJSON(e, "/v2/email/outbound-emails", `{"FromEmailAddress":"from@example.com","Content":{"Raw":{"Data":"cmF3"}}}`)

	assert.Equal(t, 200, rec.Code)
	assert.JSONEq(t, `{"MessageId":"ses-444444444"}`, rec.Body.String())
	require.NotNil(t, f2.rawIn)
	assert.Nil(t, f2.simpleIn)
	assert.Equal(t, "cmF3", f2.rawIn.Content.Raw.Data)
}

func TestAwsHandler_V2Errors(t *testing.T) {
	path := "/v2/email/outbound-emails"

	t.Run("不正な JSON は 400", func(t *testing.T) {
		rec := postJSON(newAwsServer(&fakeV1{}, &fakeV2{}), path, `{`)
		assert.Equal(t, 400, rec.Code)
	})

	t.Run("Content が空なら 400", func(t *testing.T) {
		rec := postJSON(newAwsServer(&fakeV1{}, &fakeV2{}), path, `{"Content":{}}`)
		assert.Equal(t, 400, rec.Code)
		assert.JSONEq(t, `"Invalid content"`, rec.Body.String())
	})

	t.Run("Simple の送信失敗は 500", func(t *testing.T) {
		rec := postJSON(newAwsServer(&fakeV1{}, &fakeV2{err: errors.New("boom")}), path, `{"Content":{"Simple":{}}}`)
		assert.Equal(t, 500, rec.Code)
	})

	t.Run("Raw の送信失敗は 500", func(t *testing.T) {
		rec := postJSON(newAwsServer(&fakeV1{}, &fakeV2{err: errors.New("boom")}), path, `{"Content":{"Raw":{"Data":"x"}}}`)
		assert.Equal(t, 500, rec.Code)
	})
}
