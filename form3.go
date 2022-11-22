package form3

import (
	"net/http"

	"github.com/google/uuid"
)

const (
	DefaultBaseUrl = "https://api.form3.tech"
)

// Client is the Form3 API client.
type Client struct {
	api Api
	// Accounts is the Form3 API client for /v1/organisation/accounts endpoints.
	accounts AccountsClient

	// uuidProvider returns unique UUIDv4 identifiers used as ID of new Form3 API resources.
	uuidProvider func() string
	// httpClient is an instance of http.Client used for API requests.
	httpClient *http.Client
	// backOffProvider returns a fresh BackOff instance that governs API endpoint retry policy.
	backOffProvider func() BackOff
	baseUrl         string
	organisationId  string
}

func (c *Client) Api() Api {
	return c.api
}

func (c *Client) Accounts() AccountsClient {
	return c.accounts
}

func New() *Client {
	client := &Client{
		baseUrl: DefaultBaseUrl,
		uuidProvider: func() string {
			return uuid.NewString()
		},
		backOffProvider: DefaultBackOffProvider,
		httpClient:      &http.Client{},
	}

	client.api = &api{c: client}
	client.accounts = &accountsClient{c: client}

	return client
}
