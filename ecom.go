package vipps

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Ecom is the interface that wraps all the methods available for interacting
// with the Vipps Ecom API.
type Ecom interface {
	InitiatePayment(ctx context.Context, cmd InitiatePaymentCommand) (*PaymentReference, error)
	CapturePayment(ctx context.Context, cmd CapturePaymentCommand) (*CapturedPayment, error)
	CancelPayment(ctx context.Context, cmd CancelPaymentCommand) (*CancelledPayment, error)
	RefundPayment(ctx context.Context, cmd RefundPaymentCommand) (*RefundedPayment, error)
	GetPayment(ctx context.Context, orderID string) (*Payment, error)
}

// ErrEcomm represents errors returned from the Vipps Ecom API.
type ErrEcom []EcomAPIError

// EcomAPIError represents a single error returned from the Vipps Ecom API.
type EcomAPIError struct {
	Group   string `json:"errorGroup"`
	Message string `json:"errorMessage"`
	Code    string `json:"errorCode"`
}

func (e ErrEcom) Error() string {
	s := []string{"vipps:"}
	if len(e) > 1 {
		s = append(s, "multiple errors:")
	}
	for _, e := range e {
		s = append(s, fmt.Sprintf("[%s] %s (code %s)", e.Group, e.Message, e.Code))
	}
	return strings.Join(s, " ")
}

// Timestamp is a time.Time with a custom JSON marshaller.
type Timestamp time.Time

func (t Timestamp) MarshalJSON() ([]byte, error) {
	b := fmt.Sprintf("%q", time.Time(t).Format(time.RFC3339))
	return []byte(b), nil
}

// RefundPaymentCommand represents the command used to refund Vipps Ecom
// payments
type RefundPaymentCommand struct {
	IdempotencyKey       string
	OrderID              string
	MerchantSerialNumber string
	TransactionText      string
	Amount               int
}

// RefundedPayment represents a refunded Vipps Ecom payment
type RefundedPayment struct {
	OrderID            string             `json:"orderId"`
	TransactionInfo    TransactionInfo    `json:"transactionInfo"`
	TransactionSummary TransactionSummary `json:"transactionSummary"`
}

// CapturePaymentCommand represents the command used to capture Vipps Ecom
// payments
type CapturePaymentCommand struct {
	IdempotencyKey       string
	OrderID              string
	MerchantSerialNumber string
	Amount               int
	TransactionText      string
}

// CapturedPayment represents a captured Vipps Ecom payment
type CapturedPayment struct {
	OrderID            string             `json:"orderId"`
	TransactionInfo    TransactionInfo    `json:"transactionInfo"`
	TransactionSummary TransactionSummary `json:"transactionSummary"`
}

// CancelPaymentCommand represents the command used to cancel Vipps Ecom
// payments
type CancelPaymentCommand struct {
	OrderID              string
	MerchantSerialNumber string
	TransactionText      string
}

// CancelledPayment represents a cancelled payment
type CancelledPayment struct {
	OrderID            string             `json:"orderId"`
	TransactionInfo    TransactionInfo    `json:"transactionInfo"`
	TransactionSummary TransactionSummary `json:"transactionSummary"`
}

// InitiatePaymentCommand represents the command used to initiate Vipps Ecom
// payments
type InitiatePaymentCommand struct {
	MerchantInfo MerchantInfo `json:"merchantInfo"`
	CustomerInfo CustomerInfo `json:"customerInfo"`
	Transaction  Transaction  `json:"transaction"`
}

// PaymentReference represents a reference to a Vipps payment
type PaymentReference struct {
	URL     string `json:"url"`
	OrderID string `json:"orderId"`
}

// Payment represents a Vipps payment.
type Payment struct {
	OrderID            string                `json:"orderId"`
	ShippingDetails    ShippingDetails       `json:"shippingDetails"`
	TransactionLog     []TransactionLogEntry `json:"transactionLogHistory"`
	TransactionSummary TransactionSummary    `json:"transactionSummary"`
	UserDetails        UserDetails           `json:"userDetails"`
}

// TransactionLogEntry represents the list of transactions associated with a
// payment
type TransactionLogEntry struct {
	Amount           int        `json:"amount"`
	Operation        string     `json:"operation"`
	OperationSuccess bool       `json:"operationSuccess"`
	RequestID        string     `json:"requestId"`
	Timestamp        *time.Time `json:"timeStamp"`
	TransactionID    string     `json:"transactionId"`
	TransactionText  string     `json:"transactionText"`
}

// UserDetails represents customer details from Vipps
type UserDetails struct {
	BankIDVerified string `json:"bankIdVerified"`
	DateOfBirth    string `json:"dateOfBirth"`
	Email          string `json:"email"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	MobileNumber   string `json:"mobileNumber"`
	SSN            string `json:"ssn"`
	UserID         string `json:"userId"`
}

// ShippingDetails represents details for a shipping method
type ShippingDetails struct {
	Address          Address `json:"address"`
	ShippingCost     int     `json:"shippingCost"`
	ShippingMethod   string  `json:"shippingMethod"`
	ShippingMethodID string  `json:"shippingMethodId"`
}

// Address represents a Customer's shipping address
type Address struct {
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	City         string `json:"city"`
	Country      string `json:"country"`
	PostCode     string `json:"postCode"`
}

// CustomerInfo represents information about a Vipps customer
type CustomerInfo struct {
	// MobileNumber is a Norwegian mobile phone number. 8 digits, no country
	// prefix
	MobileNumber int `json:"mobileNumber"`
}

// PaymentType is the type of payment flow to use for Vipps Ecom payments
type PaymentType string

// List of values that PaymentType can take
const (
	PaymentTypeRegular PaymentType = "eComm Regular Payment"
	PaymentTypeExpress PaymentType = "eComm Express Payment"
)

// MerchantInfo represents a merchant configuration to use for Vipps Ecom
// payments
type MerchantInfo struct {
	// AuthToken, if set, will be used as the `Authorization` header of the
	// callback request that Vipps does to `CallbackURL` for transaction
	// updates.
	AuthToken string `json:"authToken,omitempty"`
	// MerchantSerialNumber uniquely represents a sales unit in the Vipps
	// system.
	MerchantSerialNumber string `json:"merchantSerialNumber"`
	// CallbackURL is a publicly reachable HTTP endpoint that will receive
	// transaction updates from Vipps.
	CallbackURL string `json:"callbackPrefix,omitempty"`
	// ConsentRemovalURL is a publicly reachable HTTP endpoint that will receive
	// requests for removal of user consents, for GDPR compliance.
	ConsentRemovalURL string `json:"consentRemovalPrefix,omitempty"`
	// RedirectURL is the URL clients are redirected to after completing the
	// payment flow in the Vipps app.
	RedirectURL string `json:"fallBack,omitempty"`
	// PaymentType selects the type of payment flow to use. If not provided,
	// `PaymentTypeRegular` is assumed by Vipps.
	PaymentType PaymentType `json:"paymentType,omitempty"`
	// IsApp signals whether the payment flow is initiated from a mobile app or
	// website. Returns a deep link url of the form `vipps://` if set to `true`.
	IsApp bool `json:"isApp,omitempty"`
	// ShippingDetailsURL is a publicly reachable HTTP endpoint that will
	// receive requests from Vipps to calculate shipping costs.
	ShippingDetailsURL string `json:"shippingDetailsPrefix,omitempty"`
	// ShippingMethods is a static list of shipping methods to present to the
	// user in the Vipps app.
	ShippingMethods []StaticShippingMethod `json:"staticShippingDetails,omitempty"`
}

// YesNoEnum is a binary choice between `Yes` or `No`.
type YesNoEnum string

// List of values that YesNoEnum can take.
const (
	Yes YesNoEnum = "Y"
	No  YesNoEnum = "N"
)

// StaticShippingMethod represents a static shipping method represented to the
// user in the Vipps app.
type StaticShippingMethod struct {
	IsDefault        YesNoEnum `json:"isDefault"`
	Priority         int       `json:"priority"`
	ShippingCost     float64   `json:"shippingCost"`
	ShippingMethod   string    `json:"shippingMethod"`
	ShippingMethodID string    `json:"shippingMethodId"`
}

// Transaction represents details for a Vipps Ecom payment.
type Transaction struct {
	// OrderID represents the order that this payment refers to. Must be unique
	// for each sales unit.
	OrderID string `json:"orderId"`
	// Amount is the amount (in cents) to be paid for the order.
	Amount int `json:"amount"`
	// TransactionText is the short message displayed to the user about the
	// payment in the Vipps app.
	TransactionText string `json:"transactionText"`
}

// TransactionSummary represents a summary of captured, refunded and available
// amounts on a Vipps Ecom payment.
type TransactionSummary struct {
	CapturedAmount           int `json:"capturedAmount"`
	RefundedAmount           int `json:"refundedAmount"`
	RemainingAmountToCapture int `json:"remainingAmountToCapture"`
	RemainingAmountToRefund  int `json:"remainingAmountToRefund"`
	BankIdentificationNumber int `json:"bankIdentificationNumber"`
}

// TransactionInfo represents details about a Transaction in the Vipps Ecom
// system.
type TransactionInfo struct {
	Amount          int        `json:"amount"`
	Status          string     `json:"status"`
	Timestamp       *time.Time `json:"timeStamp"`
	TransactionID   string     `json:"transactionId"`
	TransactionText string     `json:"transactionText"`
}

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

	req, err := c.newRequest(ctx, method, endpoint, command)
	if err != nil {
		return nil, err
	}

	err = c.do(req, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

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

	req, err := c.newRequest(ctx, method, endpoint, command)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Request-ID", cmd.IdempotencyKey)
	err = c.do(req, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *client) GetPayment(ctx context.Context, orderID string) (*Payment, error) {
	endpoint := fmt.Sprintf("%s/%s/details", c.baseUrl+ecomEndpoint, orderID)
	method := http.MethodGet
	res := Payment{}

	req, err := c.newRequest(ctx, method, endpoint, nil)
	if err != nil {
		return nil, err
	}

	err = c.do(req, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *client) InitiatePayment(ctx context.Context, cmd InitiatePaymentCommand) (*PaymentReference, error) {
	endpoint := c.baseUrl + ecomEndpoint
	method := http.MethodPost
	res := PaymentReference{}

	req, err := c.newRequest(ctx, method, endpoint, cmd)
	if err != nil {
		return nil, err
	}

	err = c.do(req, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

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

	req, err := c.newRequest(ctx, method, endpoint, command)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Request-ID", cmd.IdempotencyKey)
	err = c.do(req, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
