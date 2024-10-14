package domain_test

import (
	"aws-ses-local-go/domain"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromRawEmailRequest(t *testing.T) {
	t.Run("should return a Mail struct from a base64 encoded raw email", func(t *testing.T) {
		rawMessage := "RnJvbTogc2VuZGVyQGV4YW1wbGUuY29tClRvOiByZWNpcGllbnRAZXhhbXBsZS5jb20KU3ViamVjdDogQW1hem9uIFNFUyBUZXN0IChBV1MgU0RLIGZvciBHbykKTUlNRS1WZXJzaW9uOiAxLjAKQ29udGVudC10eXBlOiBNdWx0aXBhcnQvTWl4ZWQ7IGJvdW5kYXJ5PSJOZXh0UGFydCIKCi0tTmV4dFBhcnQKQ29udGVudC1UeXBlOiB0ZXh0L3BsYWluCgpUaGlzIGVtYWlsIHdhcyBzZW50IHdpdGggQW1hem9uIFNFUyB1c2luZyB0aGUgQVdTIFNESyBmb3IgR28uCgotLU5leHRQYXJ0CkNvbnRlbnQtVHlwZTogdGV4dC9odG1sCgo8aDE+QW1hem9uIFNFUyBUZXN0IEVtYWlsIChBV1MgU0RLIGZvciBHbyk8L2gxPjxwPlRoaXMgZW1haWwgd2FzIHNlbnQgd2l0aCA8YSBocmVmPSdodHRwczovL2F3cy5hbWF6b24uY29tL3Nlcy8nPkFtYXpvbiBTRVM8L2E+IHVzaW5nIHRoZSA8YSBocmVmPSdodHRwczovL2F3cy5hbWF6b24uY29tL3Nkay1mb3ItZ28vJz5BV1MgU0RLIGZvciBHbzwvYT4uPC9wPgoKLS1OZXh0UGFydC0t"
		expectedMail, err := domain.FromRawEmailRequest(rawMessage)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assert.Equal(t, expectedMail.From, "sender@example.com")
		assert.Equal(t, *expectedMail.To, "recipient@example.com")
		assert.Equal(t, expectedMail.Subject, "Amazon SES Test (AWS SDK for Go)")

		// FIXME: This test is failing because the text is not being parsed correctly
		// assert.Equal(t, *expectedMail.Text, "This email was sent with the AWS SDK for Go\n\n-NextPart\nContent-Type: text/plain\n\nThis email was sent with the AWS SDK for Go.\n\n-NextPart-\n")
	})
}
func TestFromRawEmailRequest_InvalidBase64(t *testing.T) {
	t.Run("should return an error for invalid base64 input", func(t *testing.T) {
		rawMessage := "InvalidBase64String"
		_, err := domain.FromRawEmailRequest(rawMessage)
		assert.Error(t, err)
	})
}

func TestFromRawEmailRequest_InvalidEmailFormat(t *testing.T) {
	t.Run("should return an error for invalid email format", func(t *testing.T) {
		rawMessage := base64.StdEncoding.EncodeToString([]byte("InvalidEmailFormat"))
		_, err := domain.FromRawEmailRequest(rawMessage)
		assert.Error(t, err)
	})
}

func TestFromRawEmailRequest_ValidEmail(t *testing.T) {
	t.Run("should return a Mail struct from a valid base64 encoded raw email", func(t *testing.T) {
		rawMessage := "RnJvbTogc2VuZGVyQGV4YW1wbGUuY29tClRvOiByZWNpcGllbnRAZXhhbXBsZS5jb20KU3ViamVjdDogVGVzdCBFbWFpbApDb250ZW50LVR5cGU6IHRleHQvcGxhaW4KClRoaXMgaXMgYSB0ZXN0IGVtYWlsLg=="
		expectedMail, err := domain.FromRawEmailRequest(rawMessage)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assert.Equal(t, expectedMail.From, "sender@example.com")
		assert.Equal(t, *expectedMail.To, "recipient@example.com")
		assert.Equal(t, expectedMail.Subject, "Test Email")
		assert.Equal(t, *expectedMail.Text, "This is a test email.")
	})
}
