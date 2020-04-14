package ecom

import (
	"context"
	"fmt"
	"github.com/torfjor/go-vipps"
	"github.com/torfjor/go-vipps/internal"
	"io/ioutil"
	"log"
	"net/http"
)

// Client is the interface that wraps all the methods available for interacting
// with the Vipps Ecom V2 API.
type Client interface {
	InitiatePayment(ctx context.Context, cmd InitiatePaymentCommand) (*PaymentReference, error)
	CapturePayment(ctx context.Context, cmd CapturePaymentCommand) (*CapturedPayment, error)
	CancelPayment(ctx context.Context, cmd CancelPaymentCommand) (*CancelledPayment, error)
	RefundPayment(ctx context.Context, cmd RefundPaymentCommand) (*RefundedPayment, error)
	GetPayment(ctx context.Context, orderID string) (*Payment, error)
}

type client struct {
	baseUrl   string
	apiClient internal.APIClient
}

// NewClient returns a configured client that implements the Client interface.
func NewClient(config vipps.ClientConfig) Client {
	var baseUrl string
	var logger *log.Logger

	if config.HTTPClient == nil {
		panic("config.HTTPClient cannot be nil")
	}

	if config.Environment == vipps.EnvironmentTesting {
		baseUrl = vipps.BaseURLTesting
	} else {
		baseUrl = vipps.BaseURL
	}

	if config.Logger == nil {
		logger = log.New(ioutil.Discard, "", log.LstdFlags)
	} else {
		logger = config.Logger
	}

	return &client{
		baseUrl: baseUrl,
		apiClient: internal.APIClient{
			L: logger,
			C: config.HTTPClient,
		},
	}
}

// CancelPayment cancels an initiated payment. Errors for payments that are not
// in a cancellable state.
func (c *client) CancelPayment(ctx context.Context, cmd CancelPaymentCommand) (*CancelledPayment, error) {
	endpoint := fmt.Sprintf("%s/%s/cancel", c.baseUrl+ecomEndpoint, cmd.OrderID)
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

	req, err := c.apiClient.NewRequest(ctx, method, endpoint, command)
	if err != nil {
		return nil, err
	}

	err = c.apiClient.Do(req, &res)
	if err != nil {
		return nil, wrapErr(err)
	}

	return &res, nil
}

// CapturePayment captures reserved amounts on a Payment
func (c *client) CapturePayment(ctx context.Context, cmd CapturePaymentCommand) (*CapturedPayment, error) {
	endpoint := fmt.Sprintf("%s/%s/capture", c.baseUrl+ecomEndpoint, cmd.OrderID)
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

	req, err := c.apiClient.NewRequest(ctx, method, endpoint, command)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Request-ID", cmd.IdempotencyKey)
	err = c.apiClient.Do(req, &res)
	if err != nil {
		return nil, wrapErr(err)
	}

	return &res, nil
}

// GetPayment gets a Payment.
func (c *client) GetPayment(ctx context.Context, orderID string) (*Payment, error) {
	endpoint := fmt.Sprintf("%s/%s/details", c.baseUrl+ecomEndpoint, orderID)
	method := http.MethodGet
	res := Payment{}

	req, err := c.apiClient.NewRequest(ctx, method, endpoint, nil)
	if err != nil {
		return nil, err
	}

	err = c.apiClient.Do(req, &res)
	if err != nil {
		return nil, wrapErr(err)
	}

	return &res, nil
}

// InitiatePayment initiates a new Payment and returns a reference to a resource
// hosted by Vipps where the payment flow can continue.
func (c *client) InitiatePayment(ctx context.Context, cmd InitiatePaymentCommand) (*PaymentReference, error) {
	endpoint := c.baseUrl + ecomEndpoint
	method := http.MethodPost
	res := PaymentReference{}

	req, err := c.apiClient.NewRequest(ctx, method, endpoint, cmd)
	if err != nil {
		return nil, err
	}

	err = c.apiClient.Do(req, &res)
	if err != nil {
		return nil, wrapErr(err)
	}

	return &res, nil
}

// RefundPayment refunds already captured amounts on a Payment.
func (c *client) RefundPayment(ctx context.Context, cmd RefundPaymentCommand) (*RefundedPayment, error) {
	endpoint := fmt.Sprintf("%s/%s/refund", c.baseUrl+ecomEndpoint, cmd.OrderID)
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

	req, err := c.apiClient.NewRequest(ctx, method, endpoint, command)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Request-ID", cmd.IdempotencyKey)
	err = c.apiClient.Do(req, &res)
	if err != nil {
		return nil, wrapErr(err)
	}

	return &res, nil
}
