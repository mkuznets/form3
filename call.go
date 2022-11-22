package form3

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"mkuznets.com/go/form3/models"
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
