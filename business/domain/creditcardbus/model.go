package creditcardbus

import (
	"fmt"
	"time"

	"github.com/garnizeh/fingo/business/types/money"
	"github.com/garnizeh/fingo/business/types/name"
	"github.com/google/uuid"
)

// InvoiceStatus represents the status of an invoice.
type InvoiceStatus string

// Set of statuses for an invoice.
const (
	InvoiceStatusOpen   InvoiceStatus = "open"
	InvoiceStatusClosed InvoiceStatus = "closed"
	InvoiceStatusPaid   InvoiceStatus = "paid"
)

var invoiceStatuses = map[string]InvoiceStatus{
	"open":   InvoiceStatusOpen,
	"closed": InvoiceStatusClosed,
	"paid":   InvoiceStatusPaid,
}

// ParseInvoiceStatus attempts to parse the string into an InvoiceStatus.
func ParseInvoiceStatus(value string) (InvoiceStatus, error) {
	s, ok := invoiceStatuses[value]
	if !ok {
		return InvoiceStatus(""), fmt.Errorf("invalid invoice status: %s", value)
	}
	return s, nil
}

// MustParseInvoiceStatus attempts to parse the string into an InvoiceStatus
// and panics if it fails.
func MustParseInvoiceStatus(value string) InvoiceStatus {
	s, err := ParseInvoiceStatus(value)
	if err != nil {
		panic(err)
	}
	return s
}

// String returns the string representation of the status.
func (s InvoiceStatus) String() string {
	return string(s)
}

// CreditCard represents an individual credit card.
type CreditCard struct {
	CreditCardTimestamps
	CreditCardIdentity
	ID         uuid.UUID
	UserID     uuid.UUID
	Limit      money.Money
	ClosingDay int
	DueDay     int
	Enabled    bool
}

// NewCreditCard is what we require from clients when adding a credit card.
type NewCreditCard struct {
	CreditCardIdentity
	UserID     uuid.UUID
	ClosingDay int
	DueDay     int
	Limit      money.Money
}

// UpdateCreditCard defines what information may be provided to modify an
// existing credit card.
type UpdateCreditCard struct {
	Name       *name.Name
	Limit      *money.Money
	ClosingDay *int
	DueDay     *int
	Enabled    *bool
}

type CreditCardTimestamps struct {
	DateCreated time.Time
	DateUpdated time.Time
}

type CreditCardIdentity struct {
	Name           name.Name
	LastFourDigits string
}

// Invoice represents a credit card invoice for a specific month.
type Invoice struct {
	ReferenceMonth time.Time
	DueDate        time.Time
	DateCreated    time.Time
	DateUpdated    time.Time
	Status         InvoiceStatus
	TotalAmount    money.Money
	ID             uuid.UUID
	CreditCardID   uuid.UUID
}

// NewInvoice is what we require from clients when adding an invoice.
type NewInvoice struct {
	ReferenceMonth time.Time
	DueDate        time.Time
	Status         InvoiceStatus
	TotalAmount    money.Money
	CreditCardID   uuid.UUID
}

// UpdateInvoice defines what information may be provided to modify an
// existing invoice.
type UpdateInvoice struct {
	TotalAmount *money.Money
	Status      *InvoiceStatus
	DueDate     *time.Time
}
