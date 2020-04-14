# Go Vipps
Community maintained Go client library for the [Vipps](https://vipps.no) E-commerce and Recurring payments APIs. Please see Vipps' own documentation on their [Developer page](https://vipps.no/developer/).

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/torfjor/go-vipps)

## Installation

Install vipps-go with:

```sh
go get -u github.com/torfjor/go-vipps
```

Then, import it using:

``` go
import (
    "github.com/torfjor/go-vipps"
)
```

## Usage

Usage of the Vipps APIs requires a set of OAuth client credentials and an API subscription key. A configured `Client` wraps `oauth2.clientcredentials` to automatically refresh tokens when they expire.

```go
package main

import (
	"context"
	"log"
        "os"
        "github.com/torfjor/go-vipps"
	"time"
)

func main() {
	credentials := vipps.Credentials{
		ClientID:           os.Getenv("CLIENT_ID"),
		ClientSecret:       os.Getenv("CLIENT_SECRET"),
		APISubscriptionKey: os.Getenv("API_KEY"),
	}
	config := vipps.ClientConfig{
		Logger:      log.New(os.Stdout, "", log.LstdFlags),
		Environment: vipps.EnvironmentTesting,
		Credentials: credentials,
	}
	
	client := vipps.NewClient(config)
	
	mobileNumber := 97777776
	amount := 1000
	orderID := "8b84-0ad5258beb0f"
	transactionText := "A transaction"
	
	cmd := vipps.InitiatePaymentCommand{
		MerchantInfo: vipps.MerchantInfo{
			MerchantSerialNumber: "CHANGETHIS",
			CallbackURL:          "https://some.endpoint.no/callbacks",
			RedirectURL:          "https://some.endpoint.no/redirect",
			ConsentRemovalURL:    "https://some.endpoint.no/consentremoval",
			IsApp:                false,
			PaymentType:          vipps.PaymentTypeRegular,
		},
		CustomerInfo: vipps.CustomerInfo{
			MobileNumber: mobileNumber,
		},
		Transaction:  vipps.Transaction{
			Amount: amount,
			OrderID: orderID,
			TransactionText: transactionText,
		},
	}
	
	ctx, _ := context.WithTimeout(context.Background(), 5 * time.Second)
	p, err := client.InitiatePayment(ctx, cmd)
	if err != nil {
		log.Fatal(err)
	}
	// Do something with p
}
```
