package form3

import "net/http"

// SetBaseUrl configures the base URL used in the Form3 API requests.
func (c *Client) SetBaseUrl(v string) *Client {
	c.baseUrl = v
	return c
}

// SetOrganisationId configures the organisation ID used in the Form3 API requests.
func (c *Client) SetOrganisationId(v string) *Client {
	c.organisationId = v
	return c
}

func (c *Client) SetUuidProvider(v func() string) *Client {
	c.uuidProvider = v
	return c
}

func (c *Client) SetHttpClient(v *http.Client) *Client {
	c.httpClient = v
	return c
}

func (c *Client) SetBackOffProvider(v func() BackOff) *Client {
	c.backOffProvider = v
	return c
}
