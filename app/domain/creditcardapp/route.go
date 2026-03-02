// Package creditcardapp maintains the app layer api for the credit card domain.
package creditcardapp

import (
	"net/http"

	"github.com/garnizeh/fingo/app/sdk/authclient"
	"github.com/garnizeh/fingo/app/sdk/mid"
	"github.com/garnizeh/fingo/business/domain/creditcardbus"
	"github.com/garnizeh/fingo/foundation/logger"
	"github.com/garnizeh/fingo/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log           *logger.Logger
	AuthClient    authclient.Authenticator
	CreditCardBus creditcardbus.ExtBusiness
}

// Routes registers all the credit card routes.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.AuthClient)

	a := newApp(cfg.CreditCardBus)

	app.HandlerFunc(http.MethodPost, version, "/credit-cards", a.create, authen)
	app.HandlerFunc(http.MethodPut, version, "/credit-cards/{id}", a.update, authen)
	app.HandlerFunc(http.MethodDelete, version, "/credit-cards/{id}", a.delete, authen)
	app.HandlerFunc(http.MethodGet, version, "/credit-cards", a.query, authen)
	app.HandlerFunc(http.MethodGet, version, "/credit-cards/{id}", a.queryByID, authen)
}
