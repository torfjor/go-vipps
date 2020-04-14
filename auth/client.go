// Package auth provides a HTTP client suitable for speaking to Vipps APIs.
//
// Vipps APIs require that clients authorize with a client id, a client secret,
// and an API subscription key.
package auth

import (
	"context"
	"github.com/torfjor/go-vipps"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
)

const (
	tokenEndpoint = "/accessToken/get"
)

type Credentials struct {
	ClientID           string
	ClientSecret       string
	APISubscriptionKey string
}

type customTransport struct {
	config             clientcredentials.Config
	apiSubscriptionKey string
	rt                 http.RoundTripper
}

// NewClient returns a http.Client with a custom Transport that adds required
// headers for authorizing Vipps API clients. JWT tokens are automatically
// fetched and renewed upon expiry.
func NewClient(environment vipps.Environment, credentials vipps.Credentials) *http.Client {
	var baseUrl string
	if environment == vipps.EnvironmentTesting {
		baseUrl = vipps.BaseURLTesting
	} else {
		baseUrl = vipps.BaseURL
	}

	tr := &customTransport{
		config: clientcredentials.Config{
			ClientID:     credentials.ClientID,
			ClientSecret: credentials.ClientSecret,
			TokenURL:     baseUrl + tokenEndpoint,
		},
		apiSubscriptionKey: credentials.APISubscriptionKey,
		rt:                 &http.Transport{},
	}

	tokenClient := &http.Client{
		Transport: tr,
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, tokenClient)
	return tr.config.Client(ctx)
}

// RoundTrip satisfies interface http.RoundTripper
func (ct *customTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if request.URL.Path == tokenEndpoint {
		request.Header.Add("client_id", ct.config.ClientID)
		request.Header.Add("client_secret", ct.config.ClientSecret)
	}
	request.Header.Add("Ocp-Apim-Subscription-Key", ct.apiSubscriptionKey)

	return ct.rt.RoundTrip(request)
}
