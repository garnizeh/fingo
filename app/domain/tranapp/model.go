package tranapp

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"time"

	"github.com/garnizeh/fingo/app/sdk/errs"
	"github.com/garnizeh/fingo/business/domain/productbus"
	"github.com/garnizeh/fingo/business/domain/userbus"
	"github.com/garnizeh/fingo/business/types/money"
	"github.com/garnizeh/fingo/business/types/name"
	"github.com/garnizeh/fingo/business/types/password"
	"github.com/garnizeh/fingo/business/types/quantity"
	"github.com/garnizeh/fingo/business/types/role"
)

// Product represents an individual product.
type Product struct {
	ID          string  `json:"id"`
	UserID      string  `json:"userID"`
	Name        string  `json:"name"`
	DateCreated string  `json:"dateCreated"`
	DateUpdated string  `json:"dateUpdated"`
	Cost        float64 `json:"cost"`
	Quantity    int     `json:"quantity"`
}

// Encode implements the encoder interface.
func (app *Product) Encode() (data []byte, contentType string, err error) {
	data, err = json.Marshal(app)
	contentType = "application/json"
	return data, contentType, err
}

func toAppProduct(prd *productbus.Product) *Product {
	return &Product{
		ID:          prd.ID.String(),
		UserID:      prd.UserID.String(),
		Name:        prd.Name.String(),
		Cost:        prd.Cost.Value(),
		Quantity:    prd.Quantity.Value(),
		DateCreated: prd.DateCreated.Format(time.RFC3339),
		DateUpdated: prd.DateUpdated.Format(time.RFC3339),
	}
}

// =============================================================================

// NewTran represents an example of cross domain transaction at the
// application layer.
type NewTran struct {
	User    NewUser    `json:"user"`
	Product NewProduct `json:"product"`
}

// Decode implements the decoder interface.
func (app *NewTran) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// =============================================================================

// NewUser contains information needed to create a new user.
type NewUser struct {
	Name            string   `json:"name"`
	Email           string   `json:"email"`
	Department      string   `json:"department"`
	Password        string   `json:"password"`
	PasswordConfirm string   `json:"passwordConfirm"`
	Roles           []string `json:"roles"`
}

func toBusNewUser(app *NewUser) (userbus.NewUser, error) {
	var errors errs.FieldErrors

	roles, err := role.ParseMany(app.Roles)
	if err != nil {
		errors.Add("roles", err)
	}

	addr, err := mail.ParseAddress(app.Email)
	if err != nil {
		errors.Add("email", err)
	}

	nme, err := name.Parse(app.Name)
	if err != nil {
		errors.Add("name", err)
	}

	department, err := name.ParseNull(app.Department)
	if err != nil {
		errors.Add("department", err)
	}

	pass, err := password.ParseConfirm(app.Password, app.PasswordConfirm)
	if err != nil {
		errors.Add("password", err)
	}

	if len(errors) > 0 {
		return userbus.NewUser{}, fmt.Errorf("validate: %w", errors.ToError())
	}

	bus := userbus.NewUser{
		Name:       nme,
		Email:      *addr,
		Roles:      roles,
		Department: department,
		Password:   pass,
	}

	return bus, nil
}

// =============================================================================

// NewProduct is what we require from clients when adding a Product.
type NewProduct struct {
	Name     string  `json:"name"`
	Cost     float64 `json:"cost"`
	Quantity int     `json:"quantity"`
}

func toBusNewProduct(app NewProduct) (productbus.NewProduct, error) {
	var errors errs.FieldErrors

	name, err := name.Parse(app.Name)
	if err != nil {
		errors.Add("name", err)
	}

	cost, err := money.Parse(app.Cost)
	if err != nil {
		errors.Add("cost", err)
	}

	quantity, err := quantity.Parse(app.Quantity)
	if err != nil {
		errors.Add("quantity", err)
	}

	if len(errors) > 0 {
		return productbus.NewProduct{}, fmt.Errorf("validate: %w", errors.ToError())
	}

	bus := productbus.NewProduct{
		Name:     name,
		Cost:     cost,
		Quantity: quantity,
	}

	return bus, nil
}
