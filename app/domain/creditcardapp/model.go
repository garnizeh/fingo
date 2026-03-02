// Package creditcardapp maintains the app layer api for the credit card domain.
package creditcardapp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/garnizeh/fingo/app/sdk/errs"
	"github.com/garnizeh/fingo/app/sdk/mid"
	"github.com/garnizeh/fingo/business/domain/creditcardbus"
	"github.com/garnizeh/fingo/business/types/money"
	"github.com/garnizeh/fingo/business/types/name"
)

// CreditCard represents information about an individual credit card.
type CreditCard struct {
	ID             string  `json:"id"`
	UserID         string  `json:"userID"`
	Name           string  `json:"name"`
	LastFourDigits string  `json:"lastFourDigits"`
	DateCreated    string  `json:"dateCreated"`
	DateUpdated    string  `json:"dateUpdated"`
	Limit          float64 `json:"limit"`
	ClosingDay     int     `json:"closingDay"`
	DueDay         int     `json:"dueDay"`
	Enabled        bool    `json:"enabled"`
}

// Encode implements the encoder interface.
func (app *CreditCard) Encode() (data []byte, contentType string, err error) {
	data, err = json.Marshal(app)
	contentType = "application/json"
	return
}

func toAppCreditCard(cc *creditcardbus.CreditCard) *CreditCard {
	return &CreditCard{
		ID:             cc.ID.String(),
		UserID:         cc.UserID.String(),
		Name:           cc.Name.String(),
		Limit:          cc.Limit.Value(),
		ClosingDay:     cc.ClosingDay,
		DueDay:         cc.DueDay,
		LastFourDigits: cc.LastFourDigits,
		Enabled:        cc.Enabled,
		DateCreated:    cc.DateCreated.Format(time.RFC3339),
		DateUpdated:    cc.DateUpdated.Format(time.RFC3339),
	}
}

func toAppCreditCards(ccs []creditcardbus.CreditCard) []CreditCard {
	app := make([]CreditCard, len(ccs))
	for i := range ccs {
		app[i] = *toAppCreditCard(&ccs[i])
	}

	return app
}

// =============================================================================

// NewCreditCard defines the data needed to add a new credit card.
type NewCreditCard struct {
	Name           string  `json:"name"`
	LastFourDigits string  `json:"lastFourDigits"`
	Limit          float64 `json:"limit"`
	ClosingDay     int     `json:"closingDay"`
	DueDay         int     `json:"dueDay"`
}

// Decode implements the decoder interface.
func (app *NewCreditCard) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toBusNewCreditCard(ctx context.Context, app NewCreditCard) (creditcardbus.NewCreditCard, error) {
	var fieldErrors errs.FieldErrors

	userID, err := mid.GetUserID(ctx)
	if err != nil {
		fieldErrors.Add("userID", err)
	}

	name, err := name.Parse(app.Name)
	if err != nil {
		fieldErrors.Add("name", err)
	}

	limit, err := money.Parse(app.Limit)
	if err != nil {
		fieldErrors.Add("limit", err)
	}

	if len(fieldErrors) > 0 {
		return creditcardbus.NewCreditCard{}, fmt.Errorf("validate: %w", fieldErrors.ToError())
	}

	bus := creditcardbus.NewCreditCard{
		CreditCardIdentity: creditcardbus.CreditCardIdentity{
			Name:           name,
			LastFourDigits: app.LastFourDigits,
		},
		UserID:     userID,
		Limit:      limit,
		ClosingDay: app.ClosingDay,
		DueDay:     app.DueDay,
	}

	return bus, nil
}

// =============================================================================

// UpdateCreditCard defines the data needed to update a credit card.
type UpdateCreditCard struct {
	Name       *string  `json:"name"`
	Limit      *float64 `json:"limit"`
	ClosingDay *int     `json:"closingDay"`
	DueDay     *int     `json:"dueDay"`
	Enabled    *bool    `json:"enabled"`
}

// Decode implements the decoder interface.
func (app *UpdateCreditCard) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toBusUpdateCreditCard(app UpdateCreditCard) (creditcardbus.UpdateCreditCard, error) {
	var fieldErrors errs.FieldErrors

	var nme *name.Name
	if app.Name != nil {
		nm, err := name.Parse(*app.Name)
		if err != nil {
			fieldErrors.Add("name", err)
		}
		nme = &nm
	}

	var limit *money.Money
	if app.Limit != nil {
		lmt, err := money.Parse(*app.Limit)
		if err != nil {
			fieldErrors.Add("limit", err)
		}
		limit = &lmt
	}

	if len(fieldErrors) > 0 {
		return creditcardbus.UpdateCreditCard{}, fmt.Errorf("validate: %w", fieldErrors.ToError())
	}

	bus := creditcardbus.UpdateCreditCard{
		Name:       nme,
		Limit:      limit,
		ClosingDay: app.ClosingDay,
		DueDay:     app.DueDay,
		Enabled:    app.Enabled,
	}

	return bus, nil
}
