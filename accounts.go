package form3

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"mkuznets.com/go/form3/models"
)

// AccountsClient is the Form3 API client for /v1/organisation/accounts endpoints.
type AccountsClient interface {
	// Create a new bank account or register an existing bank account with Form3
	Create(ctx context.Context, attributes *models.AccountAttributes) (*models.AccountResource, error)
	// Fetch a single Account resource using the resource ID.
	Fetch(ctx context.Context, id string) (*models.AccountResource, error)
	// Delete an Account resource using the resource ID and the current version number.
	Delete(ctx context.Context, id string, version int) error
}

type accountsClient struct {
	c *Client
}

func (s *accountsClient) Create(ctx context.Context, attributes *models.AccountAttributes) (*models.AccountResource, error) {
	request := &models.AccountResource{
		Resource: models.Resource{
			ID:             s.c.uuidProvider(),
			OrganisationId: s.c.organisationId,
			Type:           "accounts",
		},
		Attributes: attributes,
	}
	response := &models.AccountResource{}

	call := &Call{
		Method:   "POST",
		Path:     "/v1/organisation/accounts",
		Request:  request,
		Response: response,
	}
	err := s.c.Api().Do(ctx, call)

	switch e := err.(type) {
	case nil:
		return response, nil
	case Error:
		if e.Type() == ErrorConflict {
			return s.Fetch(ctx, request.ID)
		}
	}

	return nil, err
}

func (s *accountsClient) Fetch(ctx context.Context, id string) (*models.AccountResource, error) {
	response := &models.AccountResource{}
	call := &Call{
		Method:   "GET",
		Path:     fmt.Sprintf("/v1/organisation/accounts/%s", id),
		Response: response,
	}
	if err := s.c.Api().Do(ctx, call); err != nil {
		return nil, err
	}
	return response, nil
}

func (s *accountsClient) Delete(ctx context.Context, id string, accountVersion int) error {
	call := &Call{
		Method:      "DELETE",
		Path:        fmt.Sprintf("/v1/organisation/accounts/%s", id),
		QueryParams: url.Values{},
	}
	call.QueryParams.Add("version", strconv.Itoa(accountVersion))
	return s.c.Api().Do(ctx, call)
}
