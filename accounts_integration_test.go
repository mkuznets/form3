package form3_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mkuznets.com/go/form3"
	"mkuznets.com/go/form3/internal/testutils"
	"mkuznets.com/go/form3/models"
)

const (
	organisationId = "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"
)

func newAccount(ctx context.Context, client *form3.Client) (*models.AccountResource, error) {
	return client.Accounts().Create(ctx, &models.AccountAttributes{
		AccountNumber: "21751823",
		BankID:        "200401",
		BankIDCode:    models.BankIDCodeGB,
		BaseCurrency:  models.CurrencyGBP,
		Bic:           "BARCGB22",
		Country:       form3.String(models.CountryGB),
		Iban:          "GB34BARC20040121751823",
		Name:          []string{"Jane Doe", "John Doe"},
		JointAccount:  form3.Bool(true),
	})
}

func newClient(resourceId string) *form3.Client {
	return form3.New().
		SetBaseUrl(testutils.BaseUrl()).
		SetOrganisationId(organisationId).
		SetUuidProvider(func() string { return resourceId })
}

func TestIntegration_accountsClient_Create(t *testing.T) {
	testutils.EnsureIntegration(t)

	t.Run("success", func(t *testing.T) {
		resourceId := uuid.NewString()
		client := newClient(resourceId)

		resp, err := newAccount(context.Background(), client)
		require.NoError(t, err)

		assert.Equal(t, resourceId, resp.ID)
		assert.Equal(t, organisationId, resp.OrganisationId)
		assert.Equal(t, "21751823", resp.Attributes.AccountNumber)
		assert.Equal(t, "200401", resp.Attributes.BankID)
		assert.Equal(t, models.BankIDCodeGB, resp.Attributes.BankIDCode)
		assert.Equal(t, models.CurrencyGBP, resp.Attributes.BaseCurrency)
		assert.Equal(t, "BARCGB22", resp.Attributes.Bic)
		assert.Equal(t, models.CountryGB, *resp.Attributes.Country)
		assert.Equal(t, "GB34BARC20040121751823", resp.Attributes.Iban)
		assert.Equal(t, []string{"Jane Doe", "John Doe"}, resp.Attributes.Name)
		assert.Equal(t, true, *resp.Attributes.JointAccount)
	})

	t.Run("success ID conflict", func(t *testing.T) {
		resourceId := uuid.NewString()
		client := newClient(resourceId)

		_, err := newAccount(context.Background(), client)
		require.NoError(t, err)

		resp, err := newAccount(context.Background(), client)
		require.NoError(t, err)

		assert.Equal(t, resourceId, resp.ID)
		assert.Equal(t, organisationId, resp.OrganisationId)
		assert.Equal(t, "GB34BARC20040121751823", resp.Attributes.Iban)
	})

	t.Run("bad request", func(t *testing.T) {
		client := newClient(uuid.NewString())

		attrs := &models.AccountAttributes{}
		_, err := client.Accounts().Create(context.Background(), attrs)
		require.ErrorAs(t, err, &form3.Error{})
		e := err.(form3.Error)
		assert.Contains(t, e.Error(), "HTTP 400: validation failure")
		assert.Equal(t, 400, e.StatusCode)
		assert.Equal(t, form3.ErrorClientError, e.Type())
	})
}

func TestIntegration_accountsClient_Fetch(t *testing.T) {
	testutils.EnsureIntegration(t)

	t.Run("success", func(t *testing.T) {
		resourceId := uuid.NewString()
		client := newClient(resourceId)
		_, err := newAccount(context.Background(), client)
		require.NoError(t, err)

		resp, err := client.Accounts().Fetch(context.Background(), resourceId)
		require.NoError(t, err)

		assert.Equal(t, resourceId, resp.ID)
		assert.Equal(t, organisationId, resp.OrganisationId)
		assert.Equal(t, "21751823", resp.Attributes.AccountNumber)
		assert.Equal(t, "200401", resp.Attributes.BankID)
		assert.Equal(t, models.BankIDCodeGB, resp.Attributes.BankIDCode)
		assert.Equal(t, models.CurrencyGBP, resp.Attributes.BaseCurrency)
		assert.Equal(t, "BARCGB22", resp.Attributes.Bic)
		assert.Equal(t, models.CountryGB, *resp.Attributes.Country)
		assert.Equal(t, "GB34BARC20040121751823", resp.Attributes.Iban)
		assert.Equal(t, []string{"Jane Doe", "John Doe"}, resp.Attributes.Name)
		assert.Equal(t, true, *resp.Attributes.JointAccount)
	})

	t.Run("not found", func(t *testing.T) {
		client := newClient("")
		_, err := client.Accounts().Fetch(context.Background(), uuid.NewString())
		require.ErrorContains(t, err, "HTTP 404")
	})

	t.Run("bad request", func(t *testing.T) {
		client := newClient("")
		_, err := client.Accounts().Fetch(context.Background(), "123")
		require.ErrorContains(t, err, "HTTP 400: id is not a valid uuid")
	})
}

func TestIntegration_accountsClient_Delete(t *testing.T) {
	testutils.EnsureIntegration(t)

	t.Run("success", func(t *testing.T) {
		resourceId := uuid.NewString()
		client := newClient(resourceId)
		_, err := newAccount(context.Background(), client)
		require.NoError(t, err)

		err = client.Accounts().Delete(context.Background(), resourceId, 0)
		require.NoError(t, err)

		_, err = client.Accounts().Fetch(context.Background(), resourceId)
		require.ErrorAs(t, err, &form3.Error{})
		assert.Equal(t, 404, err.(form3.Error).StatusCode)
	})

	t.Run("resource not found", func(t *testing.T) {
		client := newClient(uuid.NewString())
		_, err := newAccount(context.Background(), client)
		require.NoError(t, err)

		err = client.Accounts().Delete(context.Background(), uuid.NewString(), 0)
		require.ErrorAs(t, err, &form3.Error{})
		assert.Equal(t, 404, err.(form3.Error).StatusCode)
	})

	t.Run("invalid version", func(t *testing.T) {
		resourceId := uuid.NewString()
		client := newClient(resourceId)
		_, err := newAccount(context.Background(), client)
		require.NoError(t, err)

		err = client.Accounts().Delete(context.Background(), resourceId, 123)
		require.ErrorAs(t, err, &form3.Error{})
		assert.ErrorContains(t, err, "HTTP 409: invalid version")
	})
}
