// Package creditcardotel provides an extension for creditcardbus that adds
// otel tracking.
package creditcardotel

import (
	"context"

	"github.com/garnizeh/fingo/business/domain/creditcardbus"
	"github.com/garnizeh/fingo/business/sdk/order"
	"github.com/garnizeh/fingo/business/sdk/page"
	"github.com/garnizeh/fingo/business/sdk/sqldb"
	"github.com/garnizeh/fingo/foundation/otel"
	"github.com/google/uuid"
)

// Extension provides a wrapper for otel functionality around the creditcardbus.
type Extension struct {
	bus creditcardbus.ExtBusiness
}

// NewExtension constructs a new extension that wraps the creditcardbus with otel.
func NewExtension() creditcardbus.Extension {
	return func(bus creditcardbus.ExtBusiness) creditcardbus.ExtBusiness {
		return &Extension{
			bus: bus,
		}
	}
}

// NewWithTx does not apply otel.
func (ext *Extension) NewWithTx(tx sqldb.CommitRollbacker) (creditcardbus.ExtBusiness, error) {
	return ext.bus.NewWithTx(tx)
}

// Create applies otel to the credit card creation process.
func (ext *Extension) Create(ctx context.Context, actorID uuid.UUID, ncc creditcardbus.NewCreditCard) (creditcardbus.CreditCard, error) {
	ctx, span := otel.AddSpan(ctx, "business.creditcardbus.create")
	defer span.End()

	cc, err := ext.bus.Create(ctx, actorID, ncc)
	if err != nil {
		return creditcardbus.CreditCard{}, err
	}

	return cc, nil
}

// Update applies otel to the credit card update process.
func (ext *Extension) Update(ctx context.Context, actorID uuid.UUID, cc *creditcardbus.CreditCard, ucc creditcardbus.UpdateCreditCard) (creditcardbus.CreditCard, error) {
	ctx, span := otel.AddSpan(ctx, "business.creditcardbus.update")
	defer span.End()

	updatedCC, err := ext.bus.Update(ctx, actorID, cc, ucc)
	if err != nil {
		return creditcardbus.CreditCard{}, err
	}

	return updatedCC, nil
}

// Delete applies otel to the credit card deletion process.
func (ext *Extension) Delete(ctx context.Context, actorID uuid.UUID, cc *creditcardbus.CreditCard) error {
	ctx, span := otel.AddSpan(ctx, "business.creditcardbus.delete")
	defer span.End()

	if err := ext.bus.Delete(ctx, actorID, cc); err != nil {
		return err
	}

	return nil
}

// Query applies otel to the credit card query process.
func (ext *Extension) Query(ctx context.Context, actorID uuid.UUID, filter creditcardbus.QueryFilter, orderBy order.By, page page.Page) ([]creditcardbus.CreditCard, error) {
	ctx, span := otel.AddSpan(ctx, "business.creditcardbus.query")
	defer span.End()

	return ext.bus.Query(ctx, actorID, filter, orderBy, page)
}

// Count applies otel to the credit card count process.
func (ext *Extension) Count(ctx context.Context, actorID uuid.UUID, filter creditcardbus.QueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.creditcardbus.count")
	defer span.End()

	return ext.bus.Count(ctx, actorID, filter)
}

// QueryByID applies otel to the credit card query by ID process.
func (ext *Extension) QueryByID(ctx context.Context, actorID, ccID uuid.UUID) (creditcardbus.CreditCard, error) {
	ctx, span := otel.AddSpan(ctx, "business.creditcardbus.querybyid")
	defer span.End()

	return ext.bus.QueryByID(ctx, actorID, ccID)
}
