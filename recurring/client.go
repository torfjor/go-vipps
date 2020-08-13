package recurring

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/torfjor/go-vipps"
	"github.com/torfjor/go-vipps/internal"
	"net/http"
)

const (
	recurringEndpoint = "recurring/v2/agreements"
)

type Doer interface {
	Do(req *http.Request, v interface{}) error
	NewRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Request, error)
}

// Client represents an API client for the Vipps recurring payments API.
type Client struct {
	BaseURL   string
	APIClient Doer
}

// NewClient returns a configured Client.
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

// CreateCharge creates a Charge for an Agreement.
func (c *Client) CreateCharge(ctx context.Context, cmd CreateChargeCommand) (*ChargeReference, error) {
	endpoint := fmt.Sprintf("%s/%s/%s/charges", c.BaseURL, recurringEndpoint, cmd.AgreementID)
	method := http.MethodPost
	res := ChargeReference{}

	req, err := c.APIClient.NewRequest(ctx, method, endpoint, cmd)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Idempotency-Key", cmd.IdempotencyKey)

	err = c.APIClient.Do(req, &res)
	if err != nil {
		return nil, wrapErr(err)
	}

	return &res, nil
}

// CaptureCharge captures reserved amounts on a Charge.
func (c *Client) CaptureCharge(ctx context.Context, cmd CaptureChargeCommand) error {
	endpoint := fmt.Sprintf("%s/%s/%s/charges/%s/capture", c.BaseURL, recurringEndpoint, cmd.AgreementID, cmd.ChargeID)
	method := http.MethodPost

	req, err := c.APIClient.NewRequest(ctx, method, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Idempotency-Key", cmd.IdempotencyKey)

	err = c.APIClient.Do(req, nil)
	if err != nil {
		return wrapErr(err)
	}

	return nil
}

// RefundCharge refunds already captured amounts on a Charge.
func (c *Client) RefundCharge(ctx context.Context, cmd RefundChargeCommand) error {
	endpoint := fmt.Sprintf("%s/%s/%s/charges/%s/refund", c.BaseURL, recurringEndpoint, cmd.AgreementID, cmd.ChargeID)
	method := http.MethodPost

	req, err := c.APIClient.NewRequest(ctx, method, endpoint, cmd)
	if err != nil {
		return err
	}
	req.Header.Add("Idempotency-Key", cmd.IdempotencyKey)

	err = c.APIClient.Do(req, nil)
	if err != nil {
		return wrapErr(err)
	}

	return nil
}

// CancelCharge deletes a Charge. Will error for Charges that are not in a
// cancellable state.
func (c *Client) CancelCharge(ctx context.Context, cmd DeleteChargeCommand) (*Charge, error) {
	endpoint := fmt.Sprintf("%s/%s/%s/charges/%s", c.BaseURL, recurringEndpoint, cmd.AgreementID, cmd.ChargeID)
	method := http.MethodDelete
	res := Charge{}

	req, err := c.APIClient.NewRequest(ctx, method, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Idempotency-Key", cmd.IdempotencyKey)

	err = c.APIClient.Do(req, &res)
	if err != nil {
		return nil, wrapErr(err)
	}

	return &res, nil
}

// GetCharge gets a Charge associated with an Agreement.
func (c *Client) GetCharge(ctx context.Context, cmd GetChargeCommand) (*Charge, error) {
	endpoint := fmt.Sprintf("%s/%s/%s/charges/%s", c.BaseURL, recurringEndpoint, cmd.AgreementID, cmd.ChargeID)
	method := http.MethodGet
	res := Charge{}

	req, err := c.APIClient.NewRequest(ctx, method, endpoint, nil)
	if err != nil {
		return nil, err
	}

	err = c.APIClient.Do(req, &res)
	if err != nil {
		return nil, wrapErr(err)
	}

	return &res, nil
}

// ListCharges lists Charges associated with an Agreement.
func (c *Client) ListCharges(ctx context.Context, agreementID string, status ...ChargeStatus) ([]*Charge, error) {
	var query string

	if len(status) > 0 {
		query = fmt.Sprintf("?chargeStatus=%s", status[0])
	}
	endpoint := fmt.Sprintf("%s/%s/%s/charges%s", c.BaseURL, recurringEndpoint, agreementID, query)
	method := http.MethodGet
	res := make([]*Charge, 0)

	req, err := c.APIClient.NewRequest(ctx, method, endpoint, nil)
	if err != nil {
		return nil, err
	}

	err = c.APIClient.Do(req, &res)
	if err != nil {
		return nil, wrapErr(err)
	}

	return res, nil
}

// CreateAgreement creates an Agreement.
func (c *Client) CreateAgreement(ctx context.Context, cmd CreateAgreementCommand) (*AgreementReference, error) {
	endpoint := c.BaseURL + recurringEndpoint
	method := http.MethodPost
	res := AgreementReference{}

	req, err := c.APIClient.NewRequest(ctx, method, endpoint, cmd)
	if err != nil {
		return nil, err
	}

	err = c.APIClient.Do(req, &res)
	if err != nil {
		return nil, wrapErr(err)
	}

	return &res, nil
}

// UpdateAgreement updates an Agreement.
func (c *Client) UpdateAgreement(ctx context.Context, cmd UpdateAgreementCommand) (AgreementID, error) {
	endpoint := fmt.Sprintf("%s/%s/%s", c.BaseURL, recurringEndpoint, cmd.AgreementID)
	method := http.MethodPatch
	res := struct {
		AgreementID string `json:"agreementId"`
	}{}

	req, err := c.APIClient.NewRequest(ctx, method, endpoint, cmd)
	if err != nil {
		return "", err
	}

	err = c.APIClient.Do(req, &res)
	if err != nil {
		return "", wrapErr(err)
	}

	return res.AgreementID, nil
}

// ListAgreements lists Agreements for a sales unit.
func (c *Client) ListAgreements(ctx context.Context, status ...AgreementStatus) ([]*Agreement, error) {
	var query string
	if len(status) > 0 {
		query = fmt.Sprintf("?status=%s", status[0])
	}
	endpoint := c.BaseURL + recurringEndpoint + query
	method := http.MethodGet
	res := make([]*Agreement, 0)

	req, err := c.APIClient.NewRequest(ctx, method, endpoint, nil)
	if err != nil {
		return nil, err
	}

	err = c.APIClient.Do(req, &res)
	if err != nil {
		return nil, wrapErr(err)
	}

	return res, nil
}

// GetAgreement gets an Agreement.
func (c *Client) GetAgreement(ctx context.Context, agreementID string) (*Agreement, error) {
	endpoint := fmt.Sprintf("%s/%s/%s", c.BaseURL, recurringEndpoint, agreementID)
	method := http.MethodGet
	res := Agreement{}

	req, err := c.APIClient.NewRequest(ctx, method, endpoint, nil)
	if err != nil {
		return nil, err
	}

	err = c.APIClient.Do(req, &res)
	if err != nil {
		return nil, wrapErr(err)
	}

	return &res, nil
}
