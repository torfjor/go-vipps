package ecom

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/log"
	"github.com/torfjor/go-vipps"
	"github.com/torfjor/go-vipps/internal"
)

type Doer interface {
	Do(req *http.Request, v interface{}) error
	NewRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Request, error)
}

// Client represents an API client for the Vipps ecomm v2 API.
type Client struct {
	BaseURL   string
	APIClient Doer
}

// NewClient returns a configured Client
func NewClient(config vipps.ClientConfig) *Client {
	var baseUrl string
	var logger log.Logger

	if config.HTTPClient == nil {
		panic("config.HTTPClient cannot be nil")
	}

	if config.Environment == vipps.EnvironmentTesting {
		baseUrl = vipps.BaseURLTesting
	} else {
		baseUrl = vipps.BaseURL
	}

	if config.Logger == nil {
		logger = log.NewNopLogger()
	} else {
		logger = config.Logger
	}

	return &Client{
		BaseURL: baseUrl,
		APIClient: &internal.APIClient{
			L: logger,
			C: config.HTTPClient,
		},
	}
}

// CancelPayment cancels an initiated payment. Errors for payments that are not
// in a cancellable state.
func (c *Client) CancelPayment(ctx context.Context, cmd CancelPaymentCommand) (*CancelledPayment, error) {
	endpoint := fmt.Sprintf("%s/%s/%s/cancel", c.BaseURL, ecomEndpoint, cmd.OrderID)
	method := http.MethodPut
	res := CancelledPayment{}
	command := struct {
		MerchantInfo struct {
			MerchantSerialNumber string `json:"merchantSerialNumber"`
		} `json:"merchantInfo"`
		Transaction struct {
			TransactionText string `json:"transactionText"`
		} `json:"transaction"`
	}{}
	command.MerchantInfo.MerchantSerialNumber = cmd.MerchantSerialNumber
	command.Transaction.TransactionText = cmd.TransactionText

	req, err := c.APIClient.NewRequest(ctx, method, endpoint, command)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Merchant-Serial-Number", cmd.MerchantSerialNumber)

	err = c.APIClient.Do(req, &res)
	if err != nil {
		return nil, wrapErr(err)
	}

	return &res, nil
}

// CapturePayment captures reserved amounts on a Payment
func (c *Client) CapturePayment(ctx context.Context, cmd CapturePaymentCommand) (*CapturedPayment, error) {
	endpoint := fmt.Sprintf("%s/%s/%s/capture", c.BaseURL, ecomEndpoint, cmd.OrderID)
	method := http.MethodPost
	res := CapturedPayment{}
	command := struct {
		MerchantInfo struct {
			MerchantSerialNumber string `json:"merchantSerialNumber"`
		} `json:"merchantInfo"`
		Transaction struct {
			Amount          int    `json:"amount"`
			TransactionText string `json:"transactionText"`
		} `json:"transaction"`
	}{}
	command.MerchantInfo.MerchantSerialNumber = cmd.MerchantSerialNumber
	command.Transaction.Amount = cmd.Amount
	command.Transaction.TransactionText = cmd.TransactionText

	req, err := c.APIClient.NewRequest(ctx, method, endpoint, command)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Request-ID", cmd.IdempotencyKey)
	req.Header.Add("Merchant-Serial-Number", cmd.MerchantSerialNumber)
	err = c.APIClient.Do(req, &res)
	if err != nil {
		return nil, wrapErr(err)
	}

	return &res, nil
}

// GetPayment gets a Payment.
func (c *Client) GetPayment(ctx context.Context, orderID string, merchantSerialNumber string) (*Payment, error) {
	endpoint := fmt.Sprintf("%s/%s/%s/details", c.BaseURL, ecomEndpoint, orderID)
	method := http.MethodGet
	res := Payment{}

	req, err := c.APIClient.NewRequest(ctx, method, endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Merchant-Serial-Number", merchantSerialNumber)
	err = c.APIClient.Do(req, &res)
	if err != nil {
		return nil, wrapErr(err)
	}

	return &res, nil
}

// InitiatePayment initiates a new Payment and returns a reference to a resource
// hosted by Vipps where the payment flow can continue.
func (c *Client) InitiatePayment(ctx context.Context, cmd InitiatePaymentCommand) (*PaymentReference, error) {
	q := url.Values{
		"scopes": []string{"name"},
	}
	endpoint := fmt.Sprintf("%s/%s?%s", c.BaseURL, ecomEndpoint, q.Encode())
	method := http.MethodPost
	res := PaymentReference{}

	req, err := c.APIClient.NewRequest(ctx, method, endpoint, cmd)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Merchant-Serial-Number", cmd.MerchantInfo.MerchantSerialNumber)

	err = c.APIClient.Do(req, &res)
	if err != nil {
		return nil, wrapErr(err)
	}

	return &res, nil
}

// RefundPayment refunds already captured amounts on a Payment.
func (c *Client) RefundPayment(ctx context.Context, cmd RefundPaymentCommand) (*RefundedPayment, error) {
	endpoint := fmt.Sprintf("%s/%s/%s/refund", c.BaseURL, ecomEndpoint, cmd.OrderID)
	method := http.MethodPost
	res := RefundedPayment{}
	command := struct {
		MerchantInfo struct {
			MerchantSerialNumber string `json:"merchantSerialNumber"`
		} `json:"merchantInfo"`
		Transaction struct {
			Amount          int    `json:"amount"`
			TransactionText string `json:"transactionText"`
		} `json:"transaction"`
	}{}
	command.MerchantInfo.MerchantSerialNumber = cmd.MerchantSerialNumber
	command.Transaction.TransactionText = cmd.TransactionText
	command.Transaction.Amount = cmd.Amount

	req, err := c.APIClient.NewRequest(ctx, method, endpoint, command)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Request-ID", cmd.IdempotencyKey)
	req.Header.Add("Merchant-Serial-Number", cmd.MerchantSerialNumber)
	err = c.APIClient.Do(req, &res)
	if err != nil {
		return nil, wrapErr(err)
	}

	return &res, nil
}
