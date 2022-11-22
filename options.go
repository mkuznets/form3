package form3

import "net/http"

// SetBaseUrl configures the base URL of the Form3 API.
func (c *Client) SetBaseUrl(v string) *Client {
	c.baseUrl = v
	return c
}

// SetOrganisationId configures the organisation ID used in the Form3 API requests.
func (c *Client) SetOrganisationId(v string) *Client {
	c.organisationId = v
	return c
}

// SetUuidProvider configures the provider of UUIDv4 identifiers used as ID of new Form3 API resources. Should only be used for testing.
func (c *Client) SetUuidProvider(v func() string) *Client {
	c.uuidProvider = v
	return c
}

// SetHttpClient configures the http.Client used to access the API.
func (c *Client) SetHttpClient(v *http.Client) *Client {
	c.httpClient = v
	return c
}

// SetBackOffProvider configures the provider of fresh BackOff instances that govern API endpoint retry policy.
func (c *Client) SetBackOffProvider(v func() BackOff) *Client {
	c.backOffProvider = v
	return c
}
