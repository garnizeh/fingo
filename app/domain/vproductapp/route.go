package vproductapp

import (
	"net/http"

	"github.com/garnizeh/fingo/app/sdk/auth"
	"github.com/garnizeh/fingo/app/sdk/authclient"
	"github.com/garnizeh/fingo/app/sdk/mid"
	"github.com/garnizeh/fingo/business/domain/userbus"
	"github.com/garnizeh/fingo/business/domain/vproductbus"
	"github.com/garnizeh/fingo/foundation/logger"
	"github.com/garnizeh/fingo/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log         *logger.Logger
	UserBus     userbus.ExtBusiness
	VProductBus vproductbus.ExtBusiness
	AuthClient  authclient.Authenticator
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.AuthClient)
	ruleAdmin := mid.Authorize(cfg.AuthClient, auth.RuleAdminOnly)

	api := newApp(cfg.VProductBus)

	app.HandlerFunc(http.MethodGet, version, "/vproducts", api.query, authen, ruleAdmin)
}
