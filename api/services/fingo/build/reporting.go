//go:build reporting

package build

import (
	"github.com/garnizeh/fingo/app/domain/checkapp"
	"github.com/garnizeh/fingo/app/domain/vproductapp"
	"github.com/garnizeh/fingo/app/sdk/mux"
	"github.com/garnizeh/fingo/foundation/web"
)

// Routes binds the reporting routes for the fingo service.
func Routes() rpt {
	return rpt{}
}

type rpt struct{}

// Add implements the RouterAdder interface.
func (rpt) Add(app *web.App, cfg mux.Config) {
	checkapp.Routes(app, checkapp.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
	})

	vproductapp.Routes(app, vproductapp.Config{
		UserBus:     cfg.BusConfig.UserBus,
		VProductBus: cfg.BusConfig.VProductBus,
		AuthClient:  cfg.FinGoConfig.AuthClient,
	})
}
