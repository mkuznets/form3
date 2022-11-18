package form3

import (
	"mkuznets.com/go/form3/accounts"
	"mkuznets.com/go/form3/api"
)

const (
	DefaultBaseUrl = "https://api.form3.tech"
)

// Client is the Form3 API client.
type Client struct {
	Api *api.Api
	// Accounts is the Form3 API client for /v1/organisation/accounts endpoints.
	Accounts *accounts.Client

	baseUrl        string
	organisationId string
}

func New(opts ...Option) (*Client, error) {
	s := &Client{
		baseUrl: DefaultBaseUrl,
	}
	for _, opt := range opts {
		opt(s)
	}

	if s.Api == nil {
		c, err := api.New(s.baseUrl, s.organisationId)
		if err != nil {
			return nil, err
		}
		s.Api = c
	}

	s.Accounts = &accounts.Client{A: s.Api}

	return s, nil
}

type Option = func(service *Client)

func WithBaseUrl(v string) Option {
	return func(f *Client) {
		f.baseUrl = v
	}
}

func WithOrganisationId(v string) Option {
	return func(f *Client) {
		f.organisationId = v
	}
}

func String(v string) *string {
	return &v
}
