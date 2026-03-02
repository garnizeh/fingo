package vproductbus

import (
	"time"

	"github.com/garnizeh/fingo/business/types/money"
	"github.com/garnizeh/fingo/business/types/name"
	"github.com/garnizeh/fingo/business/types/quantity"
	"github.com/google/uuid"
)

// Product represents an individual product with extended information.
type Product struct {
	DateCreated time.Time
	DateUpdated time.Time
	Name        name.Name
	UserName    name.Name
	Cost        money.Money
	Quantity    quantity.Quantity
	ID          uuid.UUID
	UserID      uuid.UUID
}
