package form3

import "net/http"

type ApiError struct {
	response http.Response
}

type ApiRequestError struct {
	ApiError
	ErrorMessage string `json:"error_message"`
	ErrorCode    string `json:"error_code"`
}

type ApiPermissionError struct {
	ApiError
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
