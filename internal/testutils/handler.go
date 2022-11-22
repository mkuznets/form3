package testutils

//go:generate moq -out handler_mock.go . Handler

import "net/http"

// Handler aliases http.Handler to allow for mocking.
type Handler interface {
	http.Handler
}
