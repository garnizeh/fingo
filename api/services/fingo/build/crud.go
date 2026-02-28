//go:build crud

package build

import (
	"github.com/garnizeh/fingo/app/domain/auditapp"
	"github.com/garnizeh/fingo/app/domain/checkapp"
	"github.com/garnizeh/fingo/app/domain/homeapp"
	"github.com/garnizeh/fingo/app/domain/productapp"
	"github.com/garnizeh/fingo/app/domain/tranapp"
	"github.com/garnizeh/fingo/app/domain/userapp"
	"github.com/garnizeh/fingo/app/sdk/mux"
	"github.com/garnizeh/fingo/foundation/web"
)

// Routes binds the crud routes for the fingo service.
func Routes() crud {
	return crud{}
}

type crud struct{}

// Add implements the RouterAdder interface.
func (crud) Add(app *web.App, cfg mux.Config) {
	checkapp.Routes(app, checkapp.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
	})

	homeapp.Routes(app, homeapp.Config{
		HomeBus:    cfg.BusConfig.HomeBus,
		AuthClient: cfg.FinGoConfig.AuthClient,
	})

	productapp.Routes(app, productapp.Config{
		ProductBus: cfg.BusConfig.ProductBus,
		AuthClient: cfg.FinGoConfig.AuthClient,
	})

	tranapp.Routes(app, tranapp.Config{
		UserBus:    cfg.BusConfig.UserBus,
		ProductBus: cfg.BusConfig.ProductBus,
		Log:        cfg.Log,
		AuthClient: cfg.FinGoConfig.AuthClient,
		DB:         cfg.DB,
	})

	userapp.Routes(app, userapp.Config{
		UserBus:    cfg.BusConfig.UserBus,
		AuthClient: cfg.FinGoConfig.AuthClient,
	})

	auditapp.Routes(app, auditapp.Config{
		Log:        cfg.Log,
		AuditBus:   cfg.BusConfig.AuditBus,
		AuthClient: cfg.FinGoConfig.AuthClient,
	})
}
