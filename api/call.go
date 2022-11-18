package api

import (
	"bytes"
	"context"
	"encoding/json"
	"mkuznets.com/go/form3/models"
	"net/http"
	"net/url"
)

// Call represents a single Form3 API endpoint invocation and its response.
type Call struct {
	Method      string
	Path        string
	QueryParams url.Values
	Request     any
	Response    any
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

// NewCall creates a new Call for the given API endpoint.
func NewCall(method, path string, requestResource any, responseResource any) *Call {
	c := &Call{
		Method:   method,
		Path:     path,
		Request:  requestResource,
		Response: responseResource,
	}
	return c
}
