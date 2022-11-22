package form3

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	// DefaultBaseUrl is the default base URL for the Form3 API.
	DefaultBaseUrl = "https://api.form3.tech"

	// DefaultHttpTimeout is the timeout used in the default http.Client.
	DefaultHttpTimeout = 60 * time.Second
)

// Client is the Form3 API client.
type Client struct {
	api Api
	// Accounts is the Form3 API client for /v1/organisation/accounts endpoints.
	accounts AccountsClient

	// uuidProvider returns unique UUIDv4 identifiers used as ID of new Form3 API resources.
	uuidProvider func() string
	// httpClient is an instance of http.Client used for API requests.
	httpClient      *http.Client
	backOffProvider func() BackOff
	baseUrl         string
	organisationId  string
}

// Api returns Api to access artibrary Form3 API endpoints. Most users should use specialised clients (such as AccountsClient) instead.
func (c *Client) Api() Api {
	return c.api
}

// Accounts returns AccountsClient to access /v1/organisation/accounts endpoints.
func (c *Client) Accounts() AccountsClient {
	return c.accounts
}

// New creates a new Form3 API client.
func New() *Client {
	client := &Client{
		baseUrl: DefaultBaseUrl,
		uuidProvider: func() string {
			return uuid.NewString()
		},
		backOffProvider: DefaultBackOffProvider,
		httpClient: &http.Client{
			Timeout: DefaultHttpTimeout,
		},
	}

	client.api = &api{c: client}
	client.accounts = &accountsClient{c: client}

	return client
}
