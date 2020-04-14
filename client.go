package vipps

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

const (
	baseURL           = "https://api.vipps.no"
	baseURLTesting    = "https://apitest.vipps.no"
	tokenEndpoint     = "/accessToken/get"
	ecomEndpoint      = "/ecomm/v2/payments"
	recurringEndpoint = "/recurring/v2/agreements"
)

// ErrUnexpectedResponse represents an unexpected erroneous response from the
// Vipps APIs.
type ErrUnexpectedResponse struct {
	Body   []byte
	Status int
}

func (e ErrUnexpectedResponse) Error() string {
	return fmt.Sprintf("Unexpected response from Vipps, body: %s, status: %d", e.Body, e.Status)
}

// Environment is the Vipps environment that a Client should use.
type Environment string

// List of values that Environment can take.
const (
	EnvironmentTesting Environment = "testing"
)

var (
	isEcomm     = regexp.MustCompile("ecomm")
	isRecurring = regexp.MustCompile("recurring")
)

// Client is the interface that wraps all the methods available for interacting
// with the Vipps APIs.
type Client interface {
	Recurring
	Ecom
}

type client struct {
	logger  *log.Logger
	baseUrl string
	config  ClientConfig
	c       *http.Client
}

// Credentials represents secrets used to authorize a merchant in the Vipps
// APIs.
type Credentials struct {
	APISubscriptionKey string
	ClientID           string
	ClientSecret       string
}

// ClientConfig represents the configuration to use for a Client
type ClientConfig struct {
	// Logger, if provided, will be used for logging each request to the Vipps
	// APIs
	Logger      *log.Logger
	Environment Environment
	Credentials Credentials
}

// NewClient returns a configured Vipps Client.
//
// Client automatically fetches authorization tokens and refreshes them upon
// expiry.
func NewClient(config ClientConfig) Client {
	var baseUrl string
	var logger *log.Logger
	switch config.Environment {
	case EnvironmentTesting:
		baseUrl = baseURLTesting
	default:
		baseUrl = baseURL
	}

	credentials := clientcredentials.Config{
		ClientID:     config.Credentials.ClientID,
		ClientSecret: config.Credentials.ClientSecret,
		TokenURL:     baseUrl + tokenEndpoint,
	}

	tr := &customTransport{
		config:             credentials,
		apiSubscriptionKey: config.Credentials.APISubscriptionKey,
		rt:                 &http.Transport{},
	}

	tokenClient := &http.Client{
		Transport: tr,
	}

	if config.Logger == nil {
		logger = log.New(ioutil.Discard, "", log.LstdFlags)
	} else {
		logger = config.Logger
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, tokenClient)
	httpClient := credentials.Client(ctx)
	return &client{
		logger:  logger,
		config:  config,
		baseUrl: baseUrl,
		c:       httpClient,
	}
}

type customTransport struct {
	config             clientcredentials.Config
	apiSubscriptionKey string
	rt                 http.RoundTripper
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

func (c *client) newRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Request, error) {
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewReader(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func (c *client) do(req *http.Request, v interface{}) error {
	now := time.Now()
	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}
	c.logger.Printf("[%d] %s %s %v", resp.StatusCode, req.Method, req.URL, time.Since(now))
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		p := req.URL.Path
		switch {
		case isEcomm.MatchString(p):
			var ecomErr ErrEcom
			if err = json.Unmarshal(body, &ecomErr); err == nil {
				return ecomErr
			}
		case isRecurring.MatchString(p):
			var recurringErr ErrRecurring
			if err = json.Unmarshal(body, &recurringErr); err == nil {
				return recurringErr
			}
		}
		return ErrUnexpectedResponse{
			Body:   body,
			Status: resp.StatusCode,
		}
	}
	if v == nil {
		return nil
	}
	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}
	return nil
}
