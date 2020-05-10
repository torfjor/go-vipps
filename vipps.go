package vipps

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"net/http"
)

const (
	BaseURL        = "https://api.vipps.no"
	BaseURLTesting = "https://apitest.vipps.no"
)

// ErrUnexpectedResponse represents an unexpected erroneous response from the
// Vipps APIs.
type ErrUnexpectedResponse struct {
	Body   []byte
	Status int
}

func (e ErrUnexpectedResponse) Error() string {
	return fmt.Sprintf("unexpected response from Vipps, body: %s, status: %d", e.Body, e.Status)
}

// Environment is the Vipps environment that a Client should use.
type Environment string

// List of values that Environment can take.
const (
	EnvironmentTesting Environment = "testing"
)

// Credentials represents secrets used to authenticate and authorize a client
// to the Vipps APIs.
type Credentials struct {
	APISubscriptionKey string
	ClientID           string
	ClientSecret       string
}

// ClientConfig represents the configuration to use for a Client
type ClientConfig struct {
	Environment Environment
	Logger      log.Logger
	HTTPClient  *http.Client
}
