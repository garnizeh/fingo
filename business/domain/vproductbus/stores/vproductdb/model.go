package vproductdb

import (
	"fmt"
	"time"

	"github.com/garnizeh/fingo/business/domain/vproductbus"
	"github.com/garnizeh/fingo/business/types/money"
	"github.com/garnizeh/fingo/business/types/name"
	"github.com/garnizeh/fingo/business/types/quantity"
	"github.com/google/uuid"
)

type productDB struct {
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
	Name        string    `db:"name"`
	UserName    string    `db:"user_name"`
	Cost        float64   `db:"cost"`
	Quantity    int       `db:"quantity"`
	ID          uuid.UUID `db:"product_id"`
	UserID      uuid.UUID `db:"user_id"`
}

func toBusProduct(db *productDB) (vproductbus.Product, error) {
	userName, err := name.Parse(db.UserName)
	if err != nil {
		return vproductbus.Product{}, fmt.Errorf("parse user name: %w", err)
	}

	name, err := name.Parse(db.Name)
	if err != nil {
		return vproductbus.Product{}, fmt.Errorf("parse name: %w", err)
	}

	cost, err := money.Parse(db.Cost)
	if err != nil {
		return vproductbus.Product{}, fmt.Errorf("parse cost: %w", err)
	}

	quantity, err := quantity.Parse(db.Quantity)
	if err != nil {
		return vproductbus.Product{}, fmt.Errorf("parse quantity: %w", err)
	}

	bus := vproductbus.Product{
		ID:          db.ID,
		UserID:      db.UserID,
		Name:        name,
		Cost:        cost,
		Quantity:    quantity,
		DateCreated: db.DateCreated.In(time.Local),
		DateUpdated: db.DateUpdated.In(time.Local),
		UserName:    userName,
	}

	return bus, nil
}

func toBusProducts(dbPrds []productDB) ([]vproductbus.Product, error) {
	bus := make([]vproductbus.Product, len(dbPrds))

	for i := range dbPrds {
		var err error
		bus[i], err = toBusProduct(&dbPrds[i])
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}
