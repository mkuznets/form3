package form3

import (
	"fmt"
	"net/http"
)

// ErrorType is an enumeration of possible API error types.
type ErrorType int

const (
	ErrorClientError ErrorType = iota
	ErrorConflict
	ErrorTooManyRequests
	ErrorServerError
	ErrorUnknown
)

// Error represents a Form3 API error.
type Error struct {
	StatusCode int
	RawBody    []byte

	// Returned with HTTP 400/409
	ResponseErrorMessage string `json:"error_message"`
	ResponseErrorCode    string `json:"error_code"`

	// Returned with HTTP 403
	ResponseError            string `json:"error"`
	ResponseErrorDescription string `json:"error_description"`
}

// Type returns the ErrorType of the error.
func (e Error) Type() ErrorType {
	switch {
	case e.StatusCode == http.StatusTooManyRequests:
		return ErrorTooManyRequests
	case e.StatusCode == http.StatusConflict:
		return ErrorConflict
	case e.StatusCode/100 == 4:
		return ErrorClientError
	case e.StatusCode/100 == 5:
		return ErrorServerError
	default:
		return ErrorUnknown
	}
}

func (e Error) code() string {
	if e.ResponseError != "" {
		return e.ResponseError
	} else if e.ResponseErrorCode != "" {
		return e.ResponseErrorCode
	} else if e.StatusCode != 0 {
		return fmt.Sprintf("HTTP %d", e.StatusCode)
	}
	return ""
}

func (e Error) message() string {
	if e.ResponseErrorDescription != "" {
		return e.ResponseErrorDescription
	} else if e.ResponseErrorMessage != "" {
		return e.ResponseErrorMessage
	} else if e.StatusCode != 0 {
		return http.StatusText(e.StatusCode)
	}
	return "Unrecognised error"
}

// Error returns a string representation of the error in the form of "[<code>: ]<message>".
func (e Error) Error() string {
	code := e.code()
	msg := e.message()
	if code == "" {
		return msg
	}
	return fmt.Sprintf("%s: %s", code, msg)
}
