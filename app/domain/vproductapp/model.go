package vproductapp

import (
	"encoding/json"
	"time"

	"github.com/garnizeh/fingo/business/domain/vproductbus"
)

// Product represents information about an individual product with
// extended information.
type Product struct {
	ID          string  `json:"id"`
	UserID      string  `json:"userID"`
	Name        string  `json:"name"`
	DateCreated string  `json:"dateCreated"`
	DateUpdated string  `json:"dateUpdated"`
	UserName    string  `json:"userName"`
	Cost        float64 `json:"cost"`
	Quantity    int     `json:"quantity"`
}

// Encode implements the encoder interface.
func (app *Product) Encode() (data []byte, comtentType string, err error) {
	data, err = json.Marshal(app)
	comtentType = "application/json"
	return
}

func toAppProduct(prd *vproductbus.Product) Product {
	return Product{
		ID:          prd.ID.String(),
		UserID:      prd.UserID.String(),
		Name:        prd.Name.String(),
		Cost:        prd.Cost.Value(),
		Quantity:    prd.Quantity.Value(),
		DateCreated: prd.DateCreated.Format(time.RFC3339),
		DateUpdated: prd.DateUpdated.Format(time.RFC3339),
		UserName:    prd.UserName.String(),
	}
}

func toAppProducts(prds []vproductbus.Product) []Product {
	app := make([]Product, len(prds))
	for i := range prds {
		app[i] = toAppProduct(&prds[i])
	}

	return app
}
