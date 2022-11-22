package form3_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"mkuznets.com/go/form3"
	"mkuznets.com/go/form3/internal/testutils"
)

func testBackOff(maxRetries int) func() form3.BackOff {
	return func() form3.BackOff {
		return testutils.NewTestBackOff(maxRetries)
	}
}

func TestApi_Do(t *testing.T) {

	t.Run("base URL error", func(t *testing.T) {
		api := form3.New().SetBaseUrl("__:__").Api()
		err := api.Do(context.Background(), &form3.Call{
			Method: "GET",
			Path:   "/v1/resource",
		})
		assert.ErrorContains(t, err, "first path segment in URL cannot contain colon")
	})

	t.Run("GET", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/v1/resource", func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			w.WriteHeader(http.StatusOK)
		})
		ts := httptest.NewServer(mux)
		defer ts.Close()

		api := form3.New().SetBaseUrl(ts.URL).Api()

		err := api.Do(context.Background(), &form3.Call{
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

		api := form3.New().SetBaseUrl(ts.URL).SetBackOffProvider(testBackOff(0)).Api()

		var response struct {
			Id string `json:"id"`
		}

		err := api.Do(context.Background(), &form3.Call{
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

		api := form3.New().SetBaseUrl(ts.URL).Api()

		err := api.Do(context.Background(), &form3.Call{
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

		api := form3.New().SetBaseUrl(ts.URL).Api()
		err := api.Do(context.Background(), &form3.Call{
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

		api := form3.New().SetBaseUrl(ts.URL).SetBackOffProvider(testBackOff(10)).Api()

		err := api.Do(context.Background(), &form3.Call{Method: "POST", Path: "/v1/resource"})
		assert.NoError(t, err)
		assert.Equal(t, 4, len(handlerMock.ServeHTTPCalls()))
	})

	t.Run("retriable HTTP error, retry until limit", func(t *testing.T) {
		handlerMock := failingHandlerMock(3, http.StatusInternalServerError)
		ts := httptest.NewServer(handlerMock)
		defer ts.Close()

		api := form3.New().SetBaseUrl(ts.URL).SetBackOffProvider(testBackOff(2)).Api()

		err := api.Do(context.Background(), &form3.Call{Method: "POST", Path: "/v1/resource"})
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

		api := form3.New().SetBaseUrl(ts.URL).SetBackOffProvider(testBackOff(2)).Api()

		err := api.Do(context.Background(), &form3.Call{Method: "POST", Path: "/v1/resource"})
		assert.ErrorContains(t, err, "EOF")
		assert.Equal(t, 3, len(handlerMock.ServeHTTPCalls()))
	})

	t.Run("non-retriable HTTP error", func(t *testing.T) {
		handlerMock := failingHandlerMock(3, http.StatusBadRequest)
		ts := httptest.NewServer(handlerMock)
		defer ts.Close()

		api := form3.New().SetBaseUrl(ts.URL).SetBackOffProvider(testBackOff(2)).Api()
		err := api.Do(context.Background(), &form3.Call{Method: "POST", Path: "/v1/resource"})
		assert.ErrorContains(t, err, "HTTP 400: API error message")
		assert.Equal(t, 1, len(handlerMock.ServeHTTPCalls()))
	})
}
