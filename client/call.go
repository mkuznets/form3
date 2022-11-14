package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type Call struct {
	method      string
	path        string
	queryParams url.Values
	resource    any
}

func (c *Call) data() ([]byte, error) {
	if c.resource == nil {
		return nil, nil
	}

	return json.Marshal(struct {
		Data any `json:"data"`
	}{c.resource})
}

func (c *Call) Request(ctx context.Context, baseURL *url.URL) (*http.Request, error) {
	u, err := baseURL.Parse(c.path)
	if err != nil {
		return nil, err
	}
	u.RawQuery = c.queryParams.Encode()

	body, err := c.data()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, c.method, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	return req, nil
}

func NewCall(method, path string, opts ...CallOption) *Call {
	c := &Call{
		method:      method,
		path:        path,
		queryParams: make(url.Values),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

type CallOption = func(service *Call)

func WithQueryParam(key, value string) CallOption {
	return func(c *Call) {
		c.queryParams.Add(key, value)
	}
}

func WithResource(v any) CallOption {
	return func(c *Call) {
		c.resource = v
	}
}
