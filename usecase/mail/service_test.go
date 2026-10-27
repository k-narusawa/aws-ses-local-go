package mail_test

import (
	"aws-ses-local-go/internal/testutil"
	"aws-ses-local-go/usecase/mail"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_DeleteMail(t *testing.T) {
	repo := &testutil.FakeMailRepo{}
	require.NoError(t, mail.NewService(repo).DeleteMail("m1"))
	assert.Equal(t, []string{"m1"}, repo.DeletedIDs)
}

func TestService_DeleteMails(t *testing.T) {
	repo := &testutil.FakeMailRepo{}
	require.NoError(t, mail.NewService(repo).DeleteMails())
	assert.Equal(t, 1, repo.DeleteAllN)
}

func TestService_Errors(t *testing.T) {
	s := mail.NewService(&testutil.FakeMailRepo{Err: errors.New("db down")})
	assert.EqualError(t, s.DeleteMail("m1"), "db down")
	assert.EqualError(t, s.DeleteMails(), "db down")
}
