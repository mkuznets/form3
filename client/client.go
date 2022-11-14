package client

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	baseUrl *url.URL
	client  *http.Client
}

func NewClient(baseUrl string, opts ...Option) (*Client, error) {
	s := &Client{
		client: &http.Client{},
	}
	for _, opt := range opts {
		opt(s)
	}

	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	s.baseUrl = u

	return s, nil
}

func (s *Client) Do(ctx context.Context, call *Call, resource any) error {

	req, err := call.Request(ctx, s.baseUrl)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if resp != nil {
		defer func(b io.ReadCloser) {
			if err := b.Close(); err != nil {
				panic(err)
			}
		}(resp.Body)
	}

	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return errors.New("error")
	}

	if resource != nil {
		var data struct {
			Data any `json:"data"`
		}
		data.Data = resource

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}
	}

	return nil
}

type Option = func(service *Client)
