//go:build !crud && !reporting

// Package build manages different build options.
package build

import (
	"github.com/garnizeh/fingo/app/domain/auditapp"
	"github.com/garnizeh/fingo/app/domain/checkapp"
	"github.com/garnizeh/fingo/app/domain/homeapp"
	"github.com/garnizeh/fingo/app/domain/productapp"
	"github.com/garnizeh/fingo/app/domain/rawapp"
	"github.com/garnizeh/fingo/app/domain/tranapp"
	"github.com/garnizeh/fingo/app/domain/userapp"
	"github.com/garnizeh/fingo/app/domain/vproductapp"
	"github.com/garnizeh/fingo/app/sdk/mux"
	"github.com/garnizeh/fingo/foundation/web"
)

// Routes binds all the routes for the fingo service.
func Routes() all {
	return all{}
}

type all struct{}

// Add implements the RouterAdder interface.
func (all) Add(app *web.App, cfg mux.Config) {
	checkapp.Routes(app, checkapp.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
	})

	homeapp.Routes(app, homeapp.Config{
		Log:        cfg.Log,
		HomeBus:    cfg.BusConfig.HomeBus,
		AuthClient: cfg.FinGoConfig.AuthClient,
	})

	productapp.Routes(app, productapp.Config{
		Log:        cfg.Log,
		ProductBus: cfg.BusConfig.ProductBus,
		AuthClient: cfg.FinGoConfig.AuthClient,
	})

	rawapp.Routes(app)

	tranapp.Routes(app, tranapp.Config{
		Log:        cfg.Log,
		DB:         cfg.DB,
		UserBus:    cfg.BusConfig.UserBus,
		ProductBus: cfg.BusConfig.ProductBus,
		AuthClient: cfg.FinGoConfig.AuthClient,
	})

	userapp.Routes(app, userapp.Config{
		Log:        cfg.Log,
		UserBus:    cfg.BusConfig.UserBus,
		AuthClient: cfg.FinGoConfig.AuthClient,
	})

	auditapp.Routes(app, auditapp.Config{
		Log:        cfg.Log,
		AuditBus:   cfg.BusConfig.AuditBus,
		AuthClient: cfg.FinGoConfig.AuthClient,
	})

	vproductapp.Routes(app, vproductapp.Config{
		Log:         cfg.Log,
		UserBus:     cfg.BusConfig.UserBus,
		VProductBus: cfg.BusConfig.VProductBus,
		AuthClient:  cfg.FinGoConfig.AuthClient,
	})
}
