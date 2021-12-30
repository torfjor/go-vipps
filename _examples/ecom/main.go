package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	logkit "github.com/go-kit/kit/log"

	"github.com/torfjor/go-vipps"
	"github.com/torfjor/go-vipps/auth"
	"github.com/torfjor/go-vipps/ecom"
)

var (
	ecomClient *ecom.Client
	mi         = ecom.MerchantInfo{
		MerchantSerialNumber: "CHANGETHIS",
		CallbackURL:          "https://some.endpoint.no/callbacks",
		RedirectURL:          "https://some.endpoint.no/redirect",
		ConsentRemovalURL:    "https://some.endpoint.no/consentremoval",
		IsApp:                false,
		PaymentType:          ecom.PaymentTypeExpress,
		ShippingMethods: []ecom.StaticShippingMethod{
			{
				IsDefault:        ecom.Yes,
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
)

func main() {
	authClient := auth.NewClient(vipps.EnvironmentTesting, credentials)
	ecomClient = ecom.NewClient(vipps.ClientConfig{
		HTTPClient:  authClient,
		Logger:      logkit.NewLogfmtLogger(os.Stdout),
		Environment: vipps.EnvironmentTesting,
	})

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

	payment := getPayment(orderID)
	fmt.Printf("Retreived payment: %+v\n", payment)

}

func getPayment(orderId string) *ecom.Payment {
	p, err := ecomClient.GetPayment(context.TODO(), orderId, mi.MerchantSerialNumber)
	if err != nil {
		log.Fatal(err)
	}
	return p
}

func initiatePayment(orderID, transactionText string, amount, mobileNumber int) string {
	c := ecom.InitiatePaymentCommand{
		MerchantInfo: mi,
		CustomerInfo: ecom.CustomerInfo{
			MobileNumber: mobileNumber,
		},
		Transaction: ecom.Transaction{
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

func capturePayment(orderID, transactionText string, amount int) *ecom.CapturedPayment {
	p, err := ecomClient.CapturePayment(context.TODO(), ecom.CapturePaymentCommand{
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
