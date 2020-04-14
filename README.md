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

Usage of the Vipps APIs requires a set of OAuth client credentials, and a API subscription key. A configured `Client` wraps `oauth2.clientcredentials` to automatically refresh tokens when they expire.

```go
package main

import (
    "github.com/torfjor/go-vipps"
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
}
```

See the examples in `examples/` for more complete examples.
