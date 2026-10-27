package domain_test

import (
	"aws-ses-local-go/domain"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToDomainError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantCode   domain.ErrorCode
		wantMsg    string
		wantStatus int
	}{
		{"user not found", domain.ErrUserNotFound, domain.ErrUserNotFoundCode, "user not found.", 404},
		{"todo not found", domain.ErrToDoNotFound, domain.ErrToDoNotFoundCode, "todo not found.", 404},
		{"その他のエラー", errors.New("boom"), domain.ErrUnexpectedCode, "boom", 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain.ToDomainError(tt.err)
			assert.Equal(t, tt.wantCode, got.Code)
			assert.Equal(t, tt.wantMsg, got.Message)
			assert.Equal(t, tt.wantStatus, domain.StatusMap[got.Code])
		})
	}
}

func TestError_MarshalJSON(t *testing.T) {
	b, err := json.Marshal(domain.Error{Code: domain.ErrToDoNotFoundCode, Message: "todo not found."})
	require.NoError(t, err)
	assert.JSONEq(t, `{"code":"404-002","message":"todo not found."}`, string(b))
}

func TestError_Error(t *testing.T) {
	e := domain.Error{Code: domain.ErrUnexpectedCode, Message: "boom"}
	assert.Equal(t, "code: 2, message: boom", e.Error())
}
