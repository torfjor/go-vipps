package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/torfjor/go-vipps"
	"log"
	"os"
)

var (
	ecomClient vipps.Ecom
	mi         = vipps.MerchantInfo{
		MerchantSerialNumber: "CHANGETHIS",
		CallbackURL:          "https://some.endpoint.no/callbacks",
		RedirectURL:          "https://some.endpoint.no/redirect",
		ConsentRemovalURL:    "https://some.endpoint.no/consentremoval",
		IsApp:                false,
		PaymentType:          vipps.PaymentTypeExpress,
		ShippingMethods: []vipps.StaticShippingMethod{
			{
				IsDefault:        vipps.Yes,
				Priority:         1,
				ShippingCost:     0,
				ShippingMethod:   "Posten servicepakke",
				ShippingMethodID: "123456",
			},
		},
	}
	credentials = vipps.Credentials{
		ClientID:           os.Getenv("CLIENT_ID"),
		ClientSecret:       os.Getenv("CLIENT_SECRET"),
		APISubscriptionKey: os.Getenv("API_KEY"),
	}
	config = vipps.ClientConfig{
		Logger:      log.New(os.Stdout, "", log.LstdFlags),
		Environment: vipps.EnvironmentTesting,
		Credentials: credentials,
	}
)

func main() {
	ecomClient = vipps.NewClient(config)

	mobileNumber := 97777776
	amount := 1000
	orderID := "8b84-0ad5258beb0g"
	transactionText := "A transaction"

	redirectUrl := initiatePayment(orderID, transactionText, amount, mobileNumber)
	fmt.Printf("Open %s in your web browser and complete the transaction in the Vipps app\n", redirectUrl)
	fmt.Printf("Press any key to continue.")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadByte()
	capturedPayment := capturePayment(orderID, transactionText, amount)
	fmt.Printf("Captured payment: %+v\n", capturedPayment)
}

func initiatePayment(orderID, transactionText string, amount, mobileNumber int) string {
	c := vipps.InitiatePaymentCommand{
		MerchantInfo: mi,
		CustomerInfo: vipps.CustomerInfo{
			MobileNumber: mobileNumber,
		},
		Transaction: vipps.Transaction{
			OrderID:         orderID,
			Amount:          amount,
			TransactionText: transactionText,
		},
	}
	res, err := ecomClient.InitiatePayment(context.TODO(), c)
	if err != nil {
		log.Fatal(err)
	}
	return res.URL
}

func capturePayment(orderID, transactionText string, amount int) *vipps.CapturedPayment {
	p, err := ecomClient.CapturePayment(context.TODO(), vipps.CapturePaymentCommand{
		IdempotencyKey:       "0xdeadbeef",
		OrderID:              orderID,
		MerchantSerialNumber: mi.MerchantSerialNumber,
		Amount:               amount,
		TransactionText:      transactionText,
	})
	if err != nil {
		log.Fatal(err)
	}
	return p
}
