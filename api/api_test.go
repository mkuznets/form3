package api_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io"
	"mkuznets.com/go/form3/api"
	"mkuznets.com/go/form3/testutils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApi_New(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		a, err := api.New("https://api.form3.tech", "org-id")
		assert.NoError(t, err)
		assert.Equal(t, "https://api.form3.tech", a.BaseUrl.String())
		assert.Equal(t, "org-id", a.OrganisationId)
	})
	t.Run("url error", func(t *testing.T) {
		_, err := api.New("__:__", "org-id")
		assert.Error(t, err)
	})
}

func TestApi_Do(t *testing.T) {
	t.Run("GET", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/v1/resource", func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			w.WriteHeader(http.StatusOK)
		})
		ts := httptest.NewServer(mux)
		defer ts.Close()

		a, _ := api.New(ts.URL, "org-id")
		err := a.Do(context.Background(), &api.Call{
			Method: "GET",
			Path:   "/v1/resource",
		})
		assert.NoError(t, err)
	})

	t.Run("GET, non-JSON response", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/v1/resource", func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`<html><body>hello</body></html>`))
		})
		ts := httptest.NewServer(mux)
		defer ts.Close()

		a, _ := api.New(ts.URL, "org-id")
		a.BackOffProvider = func() api.BackOff { return testutils.NewMaxRetriesBackOff(0) }

		var response struct {
			Id string `json:"id"`
		}

		err := a.Do(context.Background(), &api.Call{
			Method:   "GET",
			Path:     "/v1/resource",
			Response: &response,
		})
		assert.ErrorContains(t, err, "invalid character")
	})

	t.Run("POST, with response", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/v1/resource", func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			_, _ = io.Copy(io.Discard, r.Body)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":{"id":"resource-id"}}`))
		})
		ts := httptest.NewServer(mux)
		defer ts.Close()

		var response struct {
			Id string `json:"id"`
		}

		a, _ := api.New(ts.URL, "org-id")
		err := a.Do(context.Background(), &api.Call{
			Method:   "POST",
			Path:     "/v1/resource",
			Response: &response,
		})
		assert.NoError(t, err)
	})

	t.Run("POST, with request, with error", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/v1/resource", func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			body, err := io.ReadAll(r.Body)
			assert.NoError(t, err)

			assert.Equal(t, `{"data":{"id":"123"}}`, string(body))
			w.WriteHeader(http.StatusOK)
		})
		ts := httptest.NewServer(mux)
		defer ts.Close()

		request := struct {
			Id string `json:"id"`
		}{Id: "123"}

		a, _ := api.New(ts.URL, "org-id")
		err := a.Do(context.Background(), &api.Call{
			Method:  "POST",
			Path:    "/v1/resource",
			Request: request,
		})
		assert.NoError(t, err)
	})

}

func failingHandlerMock(failureCount int, failureCode int) *testutils.HandlerMock {
	i := 0
	handlerMock := &testutils.HandlerMock{
		ServeHTTPFunc: func(responseWriter http.ResponseWriter, request *http.Request) {
			if i < failureCount {
				responseWriter.WriteHeader(failureCode)
				_, _ = responseWriter.Write([]byte(`{"error_message": "API error message"}`))
			} else {
				responseWriter.WriteHeader(http.StatusOK)
			}
			i++
		},
	}
	return handlerMock
}

func TestApi_DoRetry(t *testing.T) {
	t.Run("retriable HTTP error, retry until success", func(t *testing.T) {
		handlerMock := failingHandlerMock(3, http.StatusInternalServerError)
		ts := httptest.NewServer(handlerMock)
		defer ts.Close()

		a, _ := api.New(ts.URL, "org-id")
		a.BackOffProvider = func() api.BackOff {
			return testutils.NewMaxRetriesBackOff(10)
		}
		err := a.Do(context.Background(), &api.Call{Method: "POST", Path: "/v1/resource"})
		assert.NoError(t, err)
		assert.Equal(t, 4, len(handlerMock.ServeHTTPCalls()))
	})

	t.Run("retriable HTTP error, retry until limit", func(t *testing.T) {
		handlerMock := failingHandlerMock(3, http.StatusInternalServerError)
		ts := httptest.NewServer(handlerMock)
		defer ts.Close()

		a, _ := api.New(ts.URL, "org-id")
		a.BackOffProvider = func() api.BackOff {
			return testutils.NewMaxRetriesBackOff(2)
		}
		err := a.Do(context.Background(), &api.Call{Method: "POST", Path: "/v1/resource"})
		assert.ErrorContains(t, err, "HTTP 500: API error message")
		assert.Equal(t, 3, len(handlerMock.ServeHTTPCalls()))
	})

	t.Run("retriable connection error", func(t *testing.T) {
		var ts *httptest.Server
		handlerMock := &testutils.HandlerMock{
			ServeHTTPFunc: func(responseWriter http.ResponseWriter, request *http.Request) {
				// close the server to simulate connection error
				ts.CloseClientConnections()
			},
		}
		ts = httptest.NewServer(handlerMock)
		defer ts.Close()

		a, _ := api.New(ts.URL, "org-id")
		a.BackOffProvider = func() api.BackOff {
			return testutils.NewMaxRetriesBackOff(2)
		}

		err := a.Do(context.Background(), &api.Call{Method: "POST", Path: "/v1/resource"})
		assert.ErrorContains(t, err, "EOF")
		assert.Equal(t, 3, len(handlerMock.ServeHTTPCalls()))
	})

	t.Run("non-retriable HTTP error", func(t *testing.T) {
		handlerMock := failingHandlerMock(3, http.StatusBadRequest)
		ts := httptest.NewServer(handlerMock)
		defer ts.Close()

		a, _ := api.New(ts.URL, "org-id")
		a.BackOffProvider = func() api.BackOff {
			return testutils.NewMaxRetriesBackOff(10)
		}
		err := a.Do(context.Background(), &api.Call{Method: "POST", Path: "/v1/resource"})
		assert.ErrorContains(t, err, "HTTP 400: API error message")
		assert.Equal(t, 1, len(handlerMock.ServeHTTPCalls()))
	})
}
