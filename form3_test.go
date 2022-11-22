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

	v := uuid.NewString()

	srv := form3.New().
		SetOrganisationId(uuid.NewString()).
		SetBaseUrl("http://127.0.0.1:8080").
		SetUuidProvider(func() string {
			return v
		})

	r, err := srv.Accounts().Create(context.Background(), &models.AccountAttributes{
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

	srv.SetOrganisationId(uuid.NewString())

	r, err = srv.Accounts().Create(context.Background(), &models.AccountAttributes{
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
