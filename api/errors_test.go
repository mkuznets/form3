package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"mkuznets.com/go/form3/api"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  api.Error
		want string
	}{
		{
			name: "empty",
			err:  api.Error{},
			want: "Unrecognised error",
		},
		{
			name: "status_code",
			err:  api.Error{StatusCode: 404},
			want: "HTTP 404: Not Found",
		},
		{
			name: "error_code_error_message",
			err:  api.Error{StatusCode: 400, ResponseErrorCode: "bad_request", ResponseErrorMessage: "Message parsing failed"},
			want: "bad_request: Message parsing failed",
		},
		{
			name: "error_message",
			err:  api.Error{StatusCode: 400, ResponseErrorMessage: "Message parsing failed"},
			want: "HTTP 400: Message parsing failed",
		},
		{
			name: "error_error_description",
			err:  api.Error{StatusCode: 401, ResponseError: "invalid_grant", ResponseErrorDescription: "Wrong email or password"},
			want: "invalid_grant: Wrong email or password",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.err.Error())
		})
	}
}

func TestError_Type(t *testing.T) {
	tests := []struct {
		name string
		err  api.Error
		want api.ErrorType
	}{
		{
			name: "bad request",
			err:  api.Error{StatusCode: 400},
			want: api.ErrorClientError,
		},
		{
			name: "not found",
			err:  api.Error{StatusCode: 404},
			want: api.ErrorClientError,
		},
		{
			name: "forbidden",
			err:  api.Error{StatusCode: 403},
			want: api.ErrorClientError,
		},
		{
			name: "conflict",
			err:  api.Error{StatusCode: 409},
			want: api.ErrorConflict,
		},
		{
			name: "too many requests",
			err:  api.Error{StatusCode: 429},
			want: api.ErrorTooManyRequests,
		},
		{
			name: "internal server error",
			err:  api.Error{StatusCode: 500},
			want: api.ErrorServerError,
		},
		{
			name: "gateway timeout",
			err:  api.Error{StatusCode: 502},
			want: api.ErrorServerError,
		},
		{
			name: "moved permanently",
			err:  api.Error{StatusCode: 301},
			want: api.ErrorUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.err.Type())
		})
	}
}
