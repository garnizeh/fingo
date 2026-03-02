// Package creditcardaudit provides an extension for creditcardbus that adds
// auditing functionality.
package creditcardaudit

import (
	"context"

	"github.com/garnizeh/fingo/business/domain/auditbus"
	"github.com/garnizeh/fingo/business/domain/creditcardbus"
	"github.com/garnizeh/fingo/business/sdk/order"
	"github.com/garnizeh/fingo/business/sdk/page"
	"github.com/garnizeh/fingo/business/sdk/sqldb"
	"github.com/garnizeh/fingo/business/types/domain"
	"github.com/google/uuid"
)

// Extension provides a wrapper for audit functionality around the creditcardbus.
type Extension struct {
	bus      creditcardbus.ExtBusiness
	auditBus auditbus.ExtBusiness
}

// NewExtension constructs a new extension that wraps the creditcardbus with audit.
func NewExtension(auditBus auditbus.ExtBusiness) creditcardbus.Extension {
	return func(bus creditcardbus.ExtBusiness) creditcardbus.ExtBusiness {
		return &Extension{
			bus:      bus,
			auditBus: auditBus,
		}
	}
}

// NewWithTx does not apply auditing.
func (ext *Extension) NewWithTx(tx sqldb.CommitRollbacker) (creditcardbus.ExtBusiness, error) {
	return ext.bus.NewWithTx(tx)
}

// Create applies auditing to the credit card creation process.
func (ext *Extension) Create(ctx context.Context, actorID uuid.UUID, ncc creditcardbus.NewCreditCard) (creditcardbus.CreditCard, error) {
	cc, err := ext.bus.Create(ctx, actorID, ncc)
	if err != nil {
		return creditcardbus.CreditCard{}, err
	}

	na := auditbus.NewAudit{
		ObjID:     cc.ID,
		ObjDomain: domain.CreditCard,
		ObjName:   cc.Name,
		ActorID:   actorID,
		Action:    "created",
		Data:      ncc,
		Message:   "credit card created",
	}

	if _, err := ext.auditBus.Create(ctx, &na); err != nil {
		return creditcardbus.CreditCard{}, err
	}

	return cc, nil
}

// Update applies auditing to the credit card update process.
func (ext *Extension) Update(ctx context.Context, actorID uuid.UUID, cc *creditcardbus.CreditCard, ucc creditcardbus.UpdateCreditCard) (creditcardbus.CreditCard, error) {
	createdCC, err := ext.bus.Update(ctx, actorID, cc, ucc)
	if err != nil {
		return creditcardbus.CreditCard{}, err
	}

	na := auditbus.NewAudit{
		ObjID:     cc.ID,
		ObjDomain: domain.CreditCard,
		ObjName:   cc.Name,
		ActorID:   actorID,
		Action:    "updated",
		Data:      ucc,
		Message:   "credit card updated",
	}

	if _, err := ext.auditBus.Create(ctx, &na); err != nil {
		return creditcardbus.CreditCard{}, err
	}

	return createdCC, nil
}

// Delete applies auditing to the credit card deletion process.
func (ext *Extension) Delete(ctx context.Context, actorID uuid.UUID, cc *creditcardbus.CreditCard) error {
	if err := ext.bus.Delete(ctx, actorID, cc); err != nil {
		return err
	}

	na := auditbus.NewAudit{
		ObjID:     cc.ID,
		ObjDomain: domain.CreditCard,
		ObjName:   cc.Name,
		ActorID:   actorID,
		Action:    "deleted",
		Data:      nil,
		Message:   "credit card deleted",
	}

	if _, err := ext.auditBus.Create(ctx, &na); err != nil {
		return err
	}

	return nil
}

// Query applies auditing to the credit card query process.
func (ext *Extension) Query(ctx context.Context, actorID uuid.UUID, filter creditcardbus.QueryFilter, orderBy order.By, page page.Page) ([]creditcardbus.CreditCard, error) {
	return ext.bus.Query(ctx, actorID, filter, orderBy, page)
}

// Count applies auditing to the credit card count process.
func (ext *Extension) Count(ctx context.Context, actorID uuid.UUID, filter creditcardbus.QueryFilter) (int, error) {
	return ext.bus.Count(ctx, actorID, filter)
}

// QueryByID applies auditing to the credit card query by ID process.
func (ext *Extension) QueryByID(ctx context.Context, actorID, ccID uuid.UUID) (creditcardbus.CreditCard, error) {
	return ext.bus.QueryByID(ctx, actorID, ccID)
}
