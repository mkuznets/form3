package form3_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"mkuznets.com/go/form3"
	"mkuznets.com/go/form3/models"
)

func SetupIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
}

func TestIntegrationService(t *testing.T) {
	SetupIntegration(t)

	srv, err := form3.New(
		form3.WithOrganisationId(uuid.New().String()),
		form3.WithBaseUrl("http://127.0.0.1:8080"))

	v := uuid.New().String()
	srv.Api.UuidProvider = func() string {
		return v
	}

	if err != nil {
		panic(err)
	}

	r, err := srv.Accounts.Create(context.Background(), &models.AccountAttributes{
		AccountNumber: "21751823",
		BankID:        "200401",
		BankIDCode:    models.BankIDCodeGB,
		BaseCurrency:  models.CurrencyGBP,
		Bic:           "BARCGB22",
		Country:       form3.String(models.CountryGB),
		Iban:          "GB34BARC20040121751823",
		Name:          []string{"Max Kuznetsov"},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(r)

	srv.Api.OrganisationId = uuid.New().String()

	r, err = srv.Accounts.Create(context.Background(), &models.AccountAttributes{
		AccountNumber: "21751823",
		BankID:        "200401",
		BankIDCode:    models.BankIDCodeGB,
		BaseCurrency:  models.CurrencyGBP,
		Bic:           "BARCGB22",
		Country:       form3.String(models.CountryGB),
		Iban:          "GB34BARC20040121751823",
		Name:          []string{"Max Kuznetsov"},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(r)

}
