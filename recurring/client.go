package recurring

import (
	"context"
	"fmt"
	"github.com/torfjor/go-vipps"
	"github.com/torfjor/go-vipps/internal"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	recurringEndpoint = "/recurring/v2/agreements"
)

// Client is the interface that wraps all the methods available for interacting
// with the Vipps Recurring Payments API.
type Client interface {
	CreateCharge(ctx context.Context, cmd CreateChargeCommand) (*ChargeReference, error)
	CaptureCharge(ctx context.Context, cmd CaptureChargeCommand) error
	RefundCharge(ctx context.Context, cmd RefundChargeCommand) error
	CancelCharge(ctx context.Context, cmd DeleteChargeCommand) (*Charge, error)
	GetCharge(ctx context.Context, cmd GetChargeCommand) (*Charge, error)
	ListCharges(ctx context.Context, agreementID string, status ...ChargeStatus) ([]*Charge, error)
	CreateAgreement(ctx context.Context, cmd CreateAgreementCommand) (*AgreementReference, error)
	UpdateAgreement(ctx context.Context, cmd UpdateAgreementCommand) (AgreementID, error)
	ListAgreements(ctx context.Context, status ...AgreementStatus) ([]*Agreement, error)
	GetAgreement(ctx context.Context, agreementID string) (*Agreement, error)
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

// CreateCharge creates a Charge for an Agreement.
func (c *client) CreateCharge(ctx context.Context, cmd CreateChargeCommand) (*ChargeReference, error) {
	endpoint := fmt.Sprintf("%s/%s/charges", c.baseUrl+recurringEndpoint, cmd.AgreementID)
	method := http.MethodPost
	res := ChargeReference{}

	req, err := c.apiClient.NewRequest(ctx, method, endpoint, cmd)
	if err != nil {
		return nil, err
	}

	err = c.apiClient.Do(req, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// CaptureCharge captures reserved amounts on a Charge.
func (c *client) CaptureCharge(ctx context.Context, cmd CaptureChargeCommand) error {
	endpoint := fmt.Sprintf("%s/%s/charges/%s/capture", c.baseUrl+recurringEndpoint, cmd.AgreementID, cmd.ChargeID)
	method := http.MethodPost

	req, err := c.apiClient.NewRequest(ctx, method, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Idempotency-Key", cmd.IdempotencyKey)

	err = c.apiClient.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

// RefundCharge refunds already captured amounts on a Charge.
func (c *client) RefundCharge(ctx context.Context, cmd RefundChargeCommand) error {
	endpoint := fmt.Sprintf("%s/%s/charges/%s/refund", c.baseUrl+recurringEndpoint, cmd.AgreementID, cmd.ChargeID)
	method := http.MethodPost

	req, err := c.apiClient.NewRequest(ctx, method, endpoint, cmd)
	if err != nil {
		return err
	}
	req.Header.Add("Idempotency-Key", cmd.IdempotencyKey)

	err = c.apiClient.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

// CancelCharge deletes a Charge. Will error for Charges that are not in a
// cancellable state.
func (c *client) CancelCharge(ctx context.Context, cmd DeleteChargeCommand) (*Charge, error) {
	endpoint := fmt.Sprintf("%s/%s/charges/%s", c.baseUrl+recurringEndpoint, cmd.AgreementID, cmd.ChargeID)
	method := http.MethodDelete
	res := Charge{}

	req, err := c.apiClient.NewRequest(ctx, method, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Idempotency-Key", cmd.IdempotencyKey)

	err = c.apiClient.Do(req, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// GetCharge gets a Charge associated with an Agreement.
func (c *client) GetCharge(ctx context.Context, cmd GetChargeCommand) (*Charge, error) {
	endpoint := fmt.Sprintf("%s/%s/charges/%s", c.baseUrl+recurringEndpoint, cmd.AgreementID, cmd.ChargeID)
	method := http.MethodGet
	res := Charge{}

	req, err := c.apiClient.NewRequest(ctx, method, endpoint, nil)
	if err != nil {
		return nil, err
	}

	err = c.apiClient.Do(req, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// ListCharges lists Charges associated with an Agreement.
func (c *client) ListCharges(ctx context.Context, agreementID string, status ...ChargeStatus) ([]*Charge, error) {
	var query string

	if len(status) > 0 {
		query = fmt.Sprintf("?chargeStatus=%s", status[0])
	}
	endpoint := fmt.Sprintf("%s/%s/charges%s", c.baseUrl+recurringEndpoint, agreementID, query)
	method := http.MethodGet
	res := make([]*Charge, 0)

	req, err := c.apiClient.NewRequest(ctx, method, endpoint, nil)
	if err != nil {
		return nil, err
	}

	err = c.apiClient.Do(req, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// CreateAgreement creates an Agreement.
func (c *client) CreateAgreement(ctx context.Context, cmd CreateAgreementCommand) (*AgreementReference, error) {
	endpoint := c.baseUrl + recurringEndpoint
	method := http.MethodPost
	res := AgreementReference{}

	req, err := c.apiClient.NewRequest(ctx, method, endpoint, cmd)
	if err != nil {
		return nil, err
	}

	err = c.apiClient.Do(req, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// UpdateAgreement updates an Agreement.
func (c *client) UpdateAgreement(ctx context.Context, cmd UpdateAgreementCommand) (AgreementID, error) {
	endpoint := fmt.Sprintf("%s/%s", c.baseUrl+recurringEndpoint, cmd.AgreementID)
	method := http.MethodPatch
	res := struct {
		AgreementID string `json:"agreementId"`
	}{}

	req, err := c.apiClient.NewRequest(ctx, method, endpoint, cmd)
	if err != nil {
		return "", err
	}

	err = c.apiClient.Do(req, &res)
	if err != nil {
		return "", err
	}

	return res.AgreementID, nil
}

// ListAgreements lists Agreements for a sales unit.
func (c *client) ListAgreements(ctx context.Context, status ...AgreementStatus) ([]*Agreement, error) {
	var query string
	if len(status) > 0 {
		query = fmt.Sprintf("?status=%s", status[0])
	}
	endpoint := c.baseUrl + recurringEndpoint + query
	method := http.MethodGet
	res := make([]*Agreement, 0)

	req, err := c.apiClient.NewRequest(ctx, method, endpoint, nil)
	if err != nil {
		return nil, err
	}

	err = c.apiClient.Do(req, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetAgreement gets an Agreement.
func (c *client) GetAgreement(ctx context.Context, agreementID string) (*Agreement, error) {
	endpoint := fmt.Sprintf("%s/%s", c.baseUrl+recurringEndpoint, agreementID)
	method := http.MethodGet
	res := Agreement{}

	req, err := c.apiClient.NewRequest(ctx, method, endpoint, nil)
	if err != nil {
		return nil, err
	}

	err = c.apiClient.Do(req, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
