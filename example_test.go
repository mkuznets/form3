package form3_test

import (
	"context"
	"fmt"

	"mkuznets.com/go/form3"
	"mkuznets.com/go/form3/models"
)

func ExampleNew() {
	client := form3.New().
		SetBaseUrl("https://api.form3.tech").
		SetOrganisationId("9d3a8910-a748-40a3-aca2-be3d4f469c05")

	// Create a new bank account
	ba, err := client.Accounts().Create(context.Background(), &models.AccountAttributes{
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
	if err != nil {
		panic(err)
	}
	fmt.Println(ba.ID)

	// Fetch bank account
	ba, err := client.Accounts().Fetch(context.Background(), "08e96610-d4ed-4de2-9a18-fcb3017b452c")

	// Delete bank account
	err := client.Accounts().Delete(context.Background(), "08e96610-d4ed-4de2-9a18-fcb3017b452c", 2)
}
