package accounts

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"mkuznets.com/go/form3/client"
	"mkuznets.com/go/form3/models"
)

type Service struct {
	c              *client.Client
	organisationId string
}

func NewService(client *client.Client, organisationID string) *Service {
	s := &Service{
		c:              client,
		organisationId: organisationID,
	}
	return s
}

func (s *Service) Create(ctx context.Context, attributes *models.AccountAttributes) (*models.AccountResource, error) {
	call := client.NewCall("POST", "/v1/organisation/accounts",
		client.WithResource(models.AccountResource{
			Resource: models.Resource{
				ID:             uuid.New().String(),
				OrganisationID: s.organisationId,
				Type:           "accounts",
			},
			Attributes: attributes,
		}),
	)

	resource := &models.AccountResource{}
	if err := s.c.Do(ctx, call, resource); err != nil {
		return nil, err
	}

	return resource, nil
}

func (s *Service) Fetch(ctx context.Context, accountId string) (*models.AccountResource, error) {
	path := fmt.Sprintf("/v1/organisation/accounts/%s", accountId)
	call := client.NewCall("GET", path)

	resource := &models.AccountResource{}
	if err := s.c.Do(ctx, call, resource); err != nil {
		return nil, err
	}

	return resource, nil
}

func (s *Service) Delete(ctx context.Context, accountID string, accountVersion int) error {
	path := fmt.Sprintf("/v1/organisation/accounts/%s", accountID)
	call := client.NewCall("DELETE", path,
		client.WithQueryParam("version", fmt.Sprintf("%d", accountVersion)))
	return s.c.Do(ctx, call, nil)
}
