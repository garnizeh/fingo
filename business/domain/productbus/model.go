package productbus

import (
	"time"

	"github.com/garnizeh/fingo/business/types/money"
	"github.com/garnizeh/fingo/business/types/name"
	"github.com/garnizeh/fingo/business/types/quantity"
	"github.com/google/uuid"
)

// Product represents an individual Product.
type Product struct {
	DateCreated time.Time
	DateUpdated time.Time
	Name        name.Name
	Cost        money.Money
	Quantity    quantity.Quantity
	ID          uuid.UUID
	UserID      uuid.UUID
}

// NewProduct is what we require from clients when adding a Product.
type NewProduct struct {
	Name     name.Name
	Cost     money.Money
	Quantity quantity.Quantity
	UserID   uuid.UUID
}

// UpdateProduct defines what information may be provided to modify an
// existing Product. All fields are optional so clients can send just the
// fields they want changed. It uses pointer fields so we can differentiate
// between a field that was not provided and a field that was provided as
// explicitly blank. Normally we do not want to use pointers to basic types but
// we make exceptions around marshalling/unmarshalling.
type UpdateProduct struct {
	Name     *name.Name
	Cost     *money.Money
	Quantity *quantity.Quantity
}
