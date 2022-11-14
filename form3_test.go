package form3_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"mkuznets.com/go/form3"
	"mkuznets.com/go/form3/client"
	"mkuznets.com/go/form3/models"
	"mkuznets.com/go/form3/testutils"
)

func TestIntegrationService(t *testing.T) {
	testutils.SetupIntegration(t)

	form3Client, err := client.NewClient("http://127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	srv, err := form3.New(
		form3.WithOrganisationId(uuid.New().String()),
		form3.WithClient(form3Client))
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

	r, err = srv.Accounts.Fetch(context.Background(), "397ad84b-0f07-4195-b927-716478fee12c")
	if err != nil {
		panic(err)
	}

	fmt.Println(r)

}
