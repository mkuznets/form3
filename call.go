package form3

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"mkuznets.com/go/form3/models"
)

// Call represents a single Form3 API endpoint invocation and an optional structure to store the response.
type Call struct {
	// Method is the HTTP method to use for the request.
	Method string
	// Path is the path to the endpoint, relative to the base URL.
	Path string
	// QueryParams is a map of query parameters to add to the request.
	QueryParams url.Values
	// Request is the JSON-serialisable struct to send as the request body. Should be nil for endpoints without JSON request body.
	Request any
	// Response is an optional pointer to a struct to unmarshal the response body into. Should be nil for endpoints without JSON response.
	Response any
}

func (c *Call) body() ([]byte, error) {
	if c.Request == nil {
		return nil, nil
	}
	return json.Marshal(models.Body{Data: c.Request})
}

func (c *Call) httpRequest(ctx context.Context, baseURL *url.URL) (*http.Request, error) {
	u, err := baseURL.Parse(c.Path)
	if err != nil {
		return nil, err
	}
	u.RawQuery = c.QueryParams.Encode()

	body, err := c.body()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, c.Method, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}
