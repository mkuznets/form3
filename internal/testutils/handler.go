package testutils

import "net/http"

//go:generate moq -out handler_mock.go . Handler
type Handler interface {
	http.Handler
}
