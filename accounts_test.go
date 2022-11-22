package form3 // Intentionally do not use `form3_test` to mock Api.

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mkuznets.com/go/form3/models"
)

func Test_accountsClient_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		apiMock := &ApiMock{
			DoFunc: func(ctx context.Context, call *Call) error {
				assert.Equal(t, "POST", call.Method)
				assert.Equal(t, "/v1/organisation/accounts", call.Path)
				assert.IsType(t, &models.AccountResource{}, call.Response)
				require.IsType(t, &models.AccountResource{}, call.Request)

				req := call.Request.(*models.AccountResource)
				assert.Equal(t, "accounts", req.Type)
				assert.Equal(t, "c52fb94b-a795-4c77-969a-74e2364edb28", req.OrganisationId)
				assert.Equal(t, "f2037281-8242-43e6-8536-0614f0b65253", req.ID)
				assert.Equal(t, "GB34BARC20040121751823", req.Attributes.Iban)
				assert.Equal(t, "BARCGB22", req.Attributes.Bic)
				assert.Equal(t, "GB", *req.Attributes.Country)

				return nil
			},
		}
		client := New().
			SetUuidProvider(func() string { return "f2037281-8242-43e6-8536-0614f0b65253" }).
			SetOrganisationId("c52fb94b-a795-4c77-969a-74e2364edb28")
		client.api = apiMock

		_, _ = client.Accounts().Create(context.Background(), &models.AccountAttributes{
			Bic:     "BARCGB22",
			Iban:    "GB34BARC20040121751823",
			Country: String(models.CountryGB),
		})
		require.Equal(t, 1, len(apiMock.calls.Do))
	})

	t.Run("idempotent conflict", func(t *testing.T) {
		apiMock := &ApiMock{
			DoFunc: func(ctx context.Context, call *Call) error {
				switch call.Method {
				case "POST":
					return Error{StatusCode: http.StatusConflict}
				case "GET":
					assert.Equal(t, "/v1/organisation/accounts/f2037281-8242-43e6-8536-0614f0b65253", call.Path)
				}
				return nil
			},
		}
		client := New().SetUuidProvider(func() string { return "f2037281-8242-43e6-8536-0614f0b65253" })
		client.api = apiMock

		_, _ = client.Accounts().Create(context.Background(), &models.AccountAttributes{})
		require.Equal(t, 2, len(apiMock.calls.Do))
		assert.Equal(t, "POST", apiMock.calls.Do[0].Call.Method)
		assert.Equal(t, "GET", apiMock.calls.Do[1].Call.Method)
	})

	t.Run("error", func(t *testing.T) {
		apiMock := &ApiMock{
			DoFunc: func(ctx context.Context, call *Call) error {
				return Error{StatusCode: http.StatusInternalServerError}
			},
		}
		client := New()
		client.api = apiMock

		_, err := client.Accounts().Create(context.Background(), &models.AccountAttributes{})
		require.ErrorAs(t, err, &Error{})
		assert.Equal(t, ErrorServerError, err.(Error).Type())
	})
}

func Test_accountsClient_Fetch(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		apiMock := &ApiMock{
			DoFunc: func(ctx context.Context, call *Call) error {
				assert.Equal(t, "GET", call.Method)
				assert.Equal(t, "/v1/organisation/accounts/123", call.Path)
				assert.Nil(t, call.Request)
				assert.IsType(t, &models.AccountResource{}, call.Response)
				return nil
			},
		}
		client := New()
		client.api = apiMock

		_, _ = client.Accounts().Fetch(context.Background(), "123")
		require.Equal(t, 1, len(apiMock.calls.Do))
	})

	t.Run("error", func(t *testing.T) {
		apiMock := &ApiMock{
			DoFunc: func(ctx context.Context, call *Call) error {
				return Error{StatusCode: http.StatusInternalServerError}
			},
		}
		client := New()
		client.api = apiMock

		_, err := client.Accounts().Fetch(context.Background(), "123")
		require.ErrorAs(t, err, &Error{})
		assert.Equal(t, ErrorServerError, err.(Error).Type())
	})
}

func Test_accountsClient_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		apiMock := &ApiMock{
			DoFunc: func(ctx context.Context, call *Call) error {
				assert.Equal(t, "DELETE", call.Method)
				assert.Equal(t, "/v1/organisation/accounts/123", call.Path)
				assert.Equal(t, "456", call.QueryParams.Get("version"))
				assert.Nil(t, call.Request)
				assert.Nil(t, call.Response)
				return nil
			},
		}
		client := New()
		client.api = apiMock

		_ = client.Accounts().Delete(context.Background(), "123", 456)
		require.Equal(t, 1, len(apiMock.calls.Do))
	})

	t.Run("error", func(t *testing.T) {
		apiMock := &ApiMock{
			DoFunc: func(ctx context.Context, call *Call) error {
				return Error{StatusCode: http.StatusInternalServerError}
			},
		}
		client := New()
		client.api = apiMock

		err := client.Accounts().Delete(context.Background(), "123", 456)
		require.ErrorAs(t, err, &Error{})
		assert.Equal(t, ErrorServerError, err.(Error).Type())
	})
}
