package accounts

import (
	"context"
	"fmt"
	"strconv"

	"mkuznets.com/go/form3/api"
	"mkuznets.com/go/form3/models"
)

// Client is the Form3 API client for /v1/organisation/accounts endpoints.
type Client struct {
	A *api.Api
}

func (s *Client) Create(ctx context.Context, attributes *models.AccountAttributes) (*models.AccountResource, error) {
	request := &models.AccountResource{
		Resource: models.Resource{
			ID:             s.A.UuidProvider(),
			OrganisationId: s.A.OrganisationId,
			Type:           "accounts",
		},
		Attributes: attributes,
	}
	response := &models.AccountResource{}

	call := &api.Call{
		Method:   "POST",
		Path:     "/v1/organisation/accounts",
		Request:  request,
		Response: response,
	}
	err := s.A.Do(ctx, call)

	switch e := err.(type) {
	case nil:
		return response, nil
	case api.Error:
		if e.Type() == api.ErrorConflict {
			return s.Fetch(ctx, request.ID)
		}
	}

	return nil, err
}

func (s *Client) Fetch(ctx context.Context, accountId string) (*models.AccountResource, error) {
	response := &models.AccountResource{}
	call := &api.Call{
		Method:   "GET",
		Path:     fmt.Sprintf("/v1/organisation/accounts/%s", accountId),
		Response: response,
	}
	if err := s.A.Do(ctx, call); err != nil {
		return nil, err
	}
	return response, nil
}

func (s *Client) Delete(ctx context.Context, accountId string, accountVersion int) error {
	call := &api.Call{
		Method: "DELETE",
		Path:   fmt.Sprintf("/v1/organisation/accounts/%s", accountId),
	}
	call.QueryParams.Add("version", strconv.Itoa(accountVersion))
	return s.A.Do(ctx, call)
}
