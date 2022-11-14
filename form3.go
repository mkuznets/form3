package form3

import (
	"mkuznets.com/go/form3/accounts"
	"mkuznets.com/go/form3/client"
)

const (
	DefaultApiUrl = "https://api.form3.tech"
)

type Form3 struct {
	Client         *client.Client
	OrganisationId string

	Accounts *accounts.Service
}

func New(opts ...Option) (*Form3, error) {
	s := &Form3{}
	for _, opt := range opts {
		opt(s)
	}

	if s.Client == nil {
		c, err := client.NewClient(DefaultApiUrl)
		if err != nil {
			return nil, err
		}
		s.Client = c
	}

	s.Accounts = accounts.NewService(s.Client, s.OrganisationId)

	return s, nil
}

type Option = func(service *Form3)

func WithClient(v *client.Client) Option {
	return func(f *Form3) {
		f.Client = v
	}
}

func WithOrganisationId(v string) Option {
	return func(f *Form3) {
		f.OrganisationId = v
	}
}
