// Package ecom provides a Client and supporting types to consume and interact
// with the Vipps Ecom V2 API
package ecom

import (
	"fmt"
	"time"
)

const ecomEndpoint = "ecomm/v2/payments"

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
	Scope        []string     `json:"scope"`
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
	// SkipLandingPage is a flag to disable the Vipps landing page and send a
	// push notification instantly.
	SkipLandingPage bool `json:"skipLandingPage"`
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

// AddressType represents an address type
type AddressType string

// List of values that AddressType can take
const (
	AddressTypeHome     AddressType = "H"
	AddressTypeBusiness AddressType = "B"
)

// ShippingCostRequest is the request sent from Vipps to calculate shipping
// costs for an order.
type ShippingCostRequest struct {
	AddressID    int         `json:"addressId"`
	AddressLine1 string      `json:"addressLine1"`
	AddressLine2 string      `json:"addressLine2"`
	City         string      `json:"city"`
	Country      string      `json:"country"`
	PostCode     string      `json:"postCode"`
	AddressType  AddressType `json:"addressType"`
}

// ShippingCostResponse is sent in response to a ShippingCostRequest
type ShippingCostResponse struct {
	AddressID       int                    `json:"addressId"`
	OrderID         string                 `json:"orderId"`
	ShippingDetails []StaticShippingMethod `json:"shippingDetails"`
}

// TransactionUpdate represents a transaction update from Vipps
type TransactionUpdate struct {
	MerchantSerialNumber string           `json:"merchantSerialNumber"`
	OrderID              string           `json:"orderId"`
	ShippingDetails      *ShippingDetails `json:"shippingDetails"`
	TransactionInfo      *TransactionInfo `json:"transactionInfo"`
	UserDetails          *UserDetails     `json:"userDetails"`
	ErrorInfo            *EcomAPIError    `json:"errorInfo"`
}
