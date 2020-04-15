// Package recurring provides a client and supporting types to consume and interact
// with the Vipps Recurring Payments API
package recurring

import (
	"time"
)

// Currency represents the currency to use for a Vipps Ecom payment.
type Currency string

// List of values that Currency can take.
const (
	CurrencyNOK Currency = "NOK"
)

// TransactionType represents the type of capture used for a payment.
type TransactionType string

// List of values that TransactionType can take.
const (
	TransactionTypeDirectCapture  TransactionType = "DIRECT_CAPTURE"
	TransactionTypeReserveCapture TransactionType = "RESERVE_CAPTURE"
)

// ChargeInterval represents the interval that recurring payments are charged.
type ChargeInterval string

// List of values that ChargeInterval can take.
const (
	ChargeIntervalMonth ChargeInterval = "MONTH"
	ChargeIntervalWeek  ChargeInterval = "WEEK"
	ChargeIntervalDay   ChargeInterval = "DAY"
)

// ChargeType represents the type of charge used for a transaction.
type ChargeType string

// List of values that ChargeType can take.
const (
	ChargeTypeInitial   ChargeType = "INITIAL"
	ChargeTypeRecurring ChargeType = "RECURRING"
)

// Campaign represents a Vipps Recurring Payments campaign.
type Campaign struct {
	Price int        `json:"campaignPrice"`
	End   *time.Time `json:"end"`
}

// InitialCharge represents the initial charge used in a Vipps recurring
// payment.
type InitialCharge struct {
	Amount          int             `json:"amount"`
	Currency        Currency        `json:"currency"`
	Description     string          `json:"description"`
	TransactionType TransactionType `json:"transactionType"`
	OrderID         string          `json:"orderId,omitempty"`
}

// CreateAgreementCommand represents the command used to create an Agreement
type CreateAgreementCommand struct {
	Campaign            *Campaign      `json:"campaign,omitempty"`
	Currency            Currency       `json:"currency"`
	CustomerPhoneNumber string         `json:"customerPhoneNumber"`
	InitialCharge       InitialCharge  `json:"initialCharge"`
	Interval            ChargeInterval `json:"interval"`
	IntervalCount       int            `json:"intervalCount"`
	IsApp               bool           `json:"isApp"`
	AgreementURL        string         `json:"merchantAgreementUrl"`
	RedirectURL         string         `json:"merchantRedirectUrl"`
	Price               int            `json:"price"`
	ProductName         string         `json:"productName"`
	ProductDescription  string         `json:"productDescription"`
}

type AgreementID = string

// AgreementReference represents a reference to an agreement associated with a
// Vipps recurring payment.
type AgreementReference struct {
	AgreementResource string `json:"agreementResource"`
	AgreementID       string `json:"agreementId"`
	URL               string `json:"vippsConfirmationUrl"`
}

// AgreementStatus is the current status of an Agreement.
type AgreementStatus string

// List of values that AgreementStatus can take.
const (
	AgreementStatusPending AgreementStatus = "PENDING"
	AgreementStatusActive  AgreementStatus = "ACTIVE"
	AgreementStatusStopped AgreementStatus = "STOPPED"
	AgreementStatusExpired AgreementStatus = "EXPIRED"
)

// Agreement represents an agreement associated with a Vipps recurring payment.
type Agreement struct {
	Campaign           *Campaign       `json:"campaign"`
	Currency           Currency        `json:"currency"`
	ID                 string          `json:"id"`
	Interval           ChargeInterval  `json:"interval"`
	IntervalCount      int             `json:"intervalCount"`
	Price              int             `json:"price"`
	ProductName        string          `json:"productName"`
	ProductDescription string          `json:"productDescription"`
	Start              *time.Time      `json:"start"`
	End                *time.Time      `json:"end"`
	Status             AgreementStatus `json:"status"`
}

// UpdateAgreementCommand represents the command used to update an Agreement
type UpdateAgreementCommand struct {
	AgreementID        string          `json:"-"`
	Campaign           *Campaign       `json:"campaign,omitempty"`
	Price              int             `json:"price,omitempty"`
	ProductName        string          `json:"productName,omitempty"`
	ProductDescription string          `json:"productDescription,omitempty"`
	Status             AgreementStatus `json:"status,omitempty"`
}

// ChargeStatus represents the current status for a Charge.
type ChargeStatus string

// List of values that ChargeStatus can take.
const (
	ChargeStatusPending           ChargeStatus = "PENDING"
	ChargeStatusDue               ChargeStatus = "DUE"
	ChargeStatusReserved          ChargeStatus = "RESERVED"
	ChargeStatusCharged           ChargeStatus = "CHARGED"
	ChargeStatusFailed            ChargeStatus = "FAILED"
	ChargeStatusCancelled         ChargeStatus = "CANCELLED"
	ChargeStatusPartiallyRefunded ChargeStatus = "PARTIALLY_REFUNDED"
	ChargeStatusRefunded          ChargeStatus = "REFUNDED"
	ChargeStatusProcessing        ChargeStatus = "PROCESSING"
)

// Charge represents a charge associated with an Agreement.
type Charge struct {
	Amount         int          `json:"amount"`
	AmountRefunded int          `json:"amountRefunded"`
	Description    string       `json:"description"`
	Due            time.Time    `json:"due"`
	ID             string       `json:"id"`
	Status         ChargeStatus `json:"status"`
	TransactionID  string       `json:"transactionId"`
	Type           ChargeType   `json:"type"`
}

// CreateChargeCommand represents the command used to created a Charge.
type CreateChargeCommand struct {
	AgreementID string   `json:"-"`
	Amount      int      `json:"amount"`
	Currency    Currency `json:"currency,omitempty"`
	Description string   `json:"description"`
	Due         DueDate  `json:"due"`
	RetryDays   int      `json:"retryDays,omitempty"`
	OrderID     string   `json:"orderId,omitempty"`
}

// DueDate is the date at which a charge is due to be paid
type DueDate struct {
	time.Time
}

func (d DueDate) MarshalJSON() ([]byte, error) {
	layout := "2006-01-02"
	return []byte(`"` + d.Time.Format(layout) + `"`), nil
}

// ChargeReference is a reference to a Charge.
type ChargeReference struct {
	ChargeID string `json:"chargeId"`
}

// ChargeIdentifier identifies a Charge.
type ChargeIdentifier struct {
	AgreementID string `json:"-"`
	ChargeID    string `json:"-"`
}

// IdempotencyKey is used to make idempotent retries in mutating commands.
type IdempotencyKey = string

// RefundChargeCommand represents the command used to refund a Charge.
type RefundChargeCommand struct {
	ChargeIdentifier `json:"-"`
	IdempotencyKey   `json:"-"`
	Amount           int    `json:"amount"`
	Description      string `json:"description"`
}

// CaptureChargeCommand represents the command used to capture a Charge.
type CaptureChargeCommand struct {
	ChargeIdentifier
	IdempotencyKey
}

// DeleteChargeCommand represents the command used to delete a Charge.
type DeleteChargeCommand struct {
	ChargeIdentifier
	IdempotencyKey
}

// GetChargeCommand represents the command used to get a Charge.
type GetChargeCommand struct {
	ChargeIdentifier
}
