package form3_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"mkuznets.com/go/form3"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  form3.Error
		want string
	}{
		{
			name: "empty",
			err:  form3.Error{},
			want: "Unrecognised error",
		},
		{
			name: "status_code",
			err:  form3.Error{StatusCode: 404},
			want: "HTTP 404: Not Found",
		},
		{
			name: "error_code_error_message",
			err:  form3.Error{StatusCode: 400, ResponseErrorCode: "bad_request", ResponseErrorMessage: "Message parsing failed"},
			want: "bad_request: Message parsing failed",
		},
		{
			name: "error_message",
			err:  form3.Error{StatusCode: 400, ResponseErrorMessage: "Message parsing failed"},
			want: "HTTP 400: Message parsing failed",
		},
		{
			name: "error_error_description",
			err:  form3.Error{StatusCode: 401, ResponseError: "invalid_grant", ResponseErrorDescription: "Wrong email or password"},
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
		err  form3.Error
		want form3.ErrorType
	}{
		{
			name: "bad request",
			err:  form3.Error{StatusCode: 400},
			want: form3.ErrorClientError,
		},
		{
			name: "not found",
			err:  form3.Error{StatusCode: 404},
			want: form3.ErrorClientError,
		},
		{
			name: "forbidden",
			err:  form3.Error{StatusCode: 403},
			want: form3.ErrorClientError,
		},
		{
			name: "conflict",
			err:  form3.Error{StatusCode: 409},
			want: form3.ErrorConflict,
		},
		{
			name: "too many requests",
			err:  form3.Error{StatusCode: 429},
			want: form3.ErrorTooManyRequests,
		},
		{
			name: "internal server error",
			err:  form3.Error{StatusCode: 500},
			want: form3.ErrorServerError,
		},
		{
			name: "gateway timeout",
			err:  form3.Error{StatusCode: 502},
			want: form3.ErrorServerError,
		},
		{
			name: "moved permanently",
			err:  form3.Error{StatusCode: 301},
			want: form3.ErrorUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.err.Type())
		})
	}
}
