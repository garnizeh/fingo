// Package creditcardbus provides business access to credit card domain.
package creditcardbus

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/garnizeh/fingo/business/domain/userbus"
	"github.com/garnizeh/fingo/business/sdk/delegate"
	"github.com/garnizeh/fingo/business/sdk/order"
	"github.com/garnizeh/fingo/business/sdk/page"
	"github.com/garnizeh/fingo/business/sdk/sqldb"
	"github.com/garnizeh/fingo/foundation/logger"
	"github.com/google/uuid"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound  = errors.New("credit card not found")
	ErrCardLimit = errors.New("credit card limit must be positive")
)

// Storer interface declares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	Create(ctx context.Context, cc CreditCard) error
	Update(ctx context.Context, cc CreditCard) error
	Delete(ctx context.Context, cc CreditCard) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]CreditCard, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, ccID uuid.UUID) (CreditCard, error)
}

// ExtBusiness interface provides support for extensions that wrap extra functionality
// around the core business logic.
type ExtBusiness interface {
	NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error)
	Create(ctx context.Context, actorID uuid.UUID, ncc NewCreditCard) (CreditCard, error)
	Update(ctx context.Context, actorID uuid.UUID, cc CreditCard, ucc UpdateCreditCard) (CreditCard, error)
	Delete(ctx context.Context, actorID uuid.UUID, cc CreditCard) error
	Query(ctx context.Context, actorID uuid.UUID, filter QueryFilter, orderBy order.By, page page.Page) ([]CreditCard, error)
	Count(ctx context.Context, actorID uuid.UUID, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, actorID uuid.UUID, ccID uuid.UUID) (CreditCard, error)
}

// Extension is a function that wraps a new layer of business logic
// around the existing business logic.
type Extension func(ExtBusiness) ExtBusiness

// Business manages the set of APIs for credit card access.
type Business struct {
	log        *logger.Logger
	userBus    userbus.ExtBusiness
	delegate   *delegate.Delegate
	storer     Storer
	extensions []Extension
}

// NewBusiness constructs a credit card business API for use.
func NewBusiness(log *logger.Logger, userBus userbus.ExtBusiness, delegate *delegate.Delegate, storer Storer, extensions ...Extension) ExtBusiness {
	b := Business{
		log:        log,
		userBus:    userBus,
		delegate:   delegate,
		storer:     storer,
		extensions: extensions,
	}

	b.registerDelegateFunctions()

	extBus := ExtBusiness(&b)

	for i := len(extensions) - 1; i >= 0; i-- {
		ext := extensions[i]
		if ext != nil {
			extBus = ext(extBus)
		}
	}

	return extBus
}

// NewWithTx constructs a new business value that will use the
// specified transaction in any store related calls.
func (b *Business) NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error) {
	storer, err := b.storer.NewWithTx(tx)
	if err != nil {
		return nil, fmt.Errorf("storer.newwithtx: %w", err)
	}

	userBus, err := b.userBus.NewWithTx(tx)
	if err != nil {
		return nil, fmt.Errorf("userbus.newwithtx: %w", err)
	}

	return &Business{
		log:        b.log,
		userBus:    userBus,
		delegate:   b.delegate,
		storer:     storer,
		extensions: b.extensions,
	}, nil
}

// Create adds a new credit card to the system.
func (b *Business) Create(ctx context.Context, actorID uuid.UUID, ncc NewCreditCard) (CreditCard, error) {
	if ncc.Limit.Value() <= 0 {
		return CreditCard{}, ErrCardLimit
	}

	now := time.Now()

	cc := CreditCard{
		ID:             uuid.New(),
		UserID:         ncc.UserID,
		Name:           ncc.Name,
		Limit:          ncc.Limit,
		ClosingDay:     ncc.ClosingDay,
		DueDay:         ncc.DueDay,
		LastFourDigits: ncc.LastFourDigits,
		Enabled:        true,
		DateCreated:    now,
		DateUpdated:    now,
	}

	if err := b.storer.Create(ctx, cc); err != nil {
		return CreditCard{}, fmt.Errorf("create: %w", err)
	}

	return cc, nil
}

// Update modifies information about a credit card.
func (b *Business) Update(ctx context.Context, actorID uuid.UUID, cc CreditCard, ucc UpdateCreditCard) (CreditCard, error) {
	if ucc.Limit != nil && ucc.Limit.Value() <= 0 {
		return CreditCard{}, ErrCardLimit
	}

	if ucc.Name != nil {
		cc.Name = *ucc.Name
	}

	if ucc.Limit != nil {
		cc.Limit = *ucc.Limit
	}

	if ucc.ClosingDay != nil {
		cc.ClosingDay = *ucc.ClosingDay
	}

	if ucc.DueDay != nil {
		cc.DueDay = *ucc.DueDay
	}

	if ucc.Enabled != nil {
		cc.Enabled = *ucc.Enabled
	}

	cc.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, cc); err != nil {
		return CreditCard{}, fmt.Errorf("update: %w", err)
	}

	return cc, nil
}

// Delete removes a credit card from the system.
func (b *Business) Delete(ctx context.Context, actorID uuid.UUID, cc CreditCard) error {
	if err := b.storer.Delete(ctx, cc); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of existing credit cards.
func (b *Business) Query(ctx context.Context, actorID uuid.UUID, filter QueryFilter, orderBy order.By, page page.Page) ([]CreditCard, error) {
	cc, err := b.storer.Query(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return cc, nil
}

// Count returns the total number of credit cards in the system.
func (b *Business) Count(ctx context.Context, actorID uuid.UUID, filter QueryFilter) (int, error) {
	count, err := b.storer.Count(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("count: %w", err)
	}

	return count, nil
}

// QueryByID finds the credit card identified by a given ID.
func (b *Business) QueryByID(ctx context.Context, actorID uuid.UUID, ccID uuid.UUID) (CreditCard, error) {
	cc, err := b.storer.QueryByID(ctx, ccID)
	if err != nil {
		return CreditCard{}, fmt.Errorf("querybyid: %w", err)
	}

	return cc, nil
}
