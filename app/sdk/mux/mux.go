// Package mux provides support to bind domain level routes
// to the application mux.
package mux

import (
	"context"
	"embed"
	"net/http"

	"github.com/garnizeh/fingo/app/sdk/auth"
	"github.com/garnizeh/fingo/app/sdk/authclient"
	"github.com/garnizeh/fingo/app/sdk/mid"
	"github.com/garnizeh/fingo/business/domain/auditbus"
	"github.com/garnizeh/fingo/business/domain/creditcardbus"
	"github.com/garnizeh/fingo/business/domain/homebus"
	"github.com/garnizeh/fingo/business/domain/productbus"
	"github.com/garnizeh/fingo/business/domain/userbus"
	"github.com/garnizeh/fingo/business/domain/vproductbus"
	"github.com/garnizeh/fingo/foundation/logger"
	"github.com/garnizeh/fingo/foundation/web"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
)

// StaticSite represents a static site to run.
type StaticSite struct {
	static     embed.FS
	staticDir  string
	staticPath string
	react      bool
}

// Options represent optional parameters.
type Options struct {
	corsOrigin []string
	sites      []StaticSite
}

// WithCORS provides configuration options for CORS.
func WithCORS(origins []string) func(opts *Options) {
	return func(opts *Options) {
		opts.corsOrigin = origins
	}
}

// WithFileServer provides configuration options for file server.
func WithFileServer(react bool, static embed.FS, dir, path string) func(opts *Options) {
	return func(opts *Options) {
		opts.sites = append(opts.sites, StaticSite{
			react:      react,
			static:     static,
			staticDir:  dir,
			staticPath: path,
		})
	}
}

// FinGoConfig contains sales service specific config.
type FinGoConfig struct {
	AuthClient authclient.Authenticator
}

// AuthConfig contains auth service specific config.
type AuthConfig struct {
	Auth *auth.Auth
}

type BusConfig struct {
	AuditBus      auditbus.ExtBusiness
	UserBus       userbus.ExtBusiness
	CreditCardBus creditcardbus.ExtBusiness
	ProductBus    productbus.ExtBusiness
	HomeBus       homebus.ExtBusiness
	VProductBus   vproductbus.ExtBusiness
}

// Config contains all the mandatory systems required by handlers.
type Config struct {
	BusConfig   BusConfig
	Tracer      trace.Tracer
	FinGoConfig FinGoConfig
	Log         *logger.Logger
	DB          *sqlx.DB
	AuthConfig  AuthConfig
	Build       string
}

// RouteAdder defines behavior that sets the routes to bind for an instance
// of the service.
type RouteAdder interface {
	Add(app *web.App, cfg *Config)
}

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI(cfg *Config, routeAdder RouteAdder, options ...func(opts *Options)) http.Handler {
	app := web.NewApp(
		cfg.Log.Info,
		cfg.Tracer,
		mid.Otel(cfg.Tracer),
		mid.Logger(cfg.Log),
		mid.Errors(cfg.Log),
		mid.Metrics(),
		mid.Panics(),
	)

	var opts Options
	for _, option := range options {
		option(&opts)
	}

	if len(opts.corsOrigin) > 0 {
		app.EnableCORS(opts.corsOrigin)
	}

	routeAdder.Add(app, cfg)

	for _, site := range opts.sites {
		switch site.react {
		case true:
			if err := app.FileServerReact(site.static, site.staticDir, site.staticPath); err != nil {
				logFileServerError(cfg, err)
			}

		default:
			if err := app.FileServer(site.static, site.staticDir, site.staticPath); err != nil {
				logFileServerError(cfg, err)
			}
		}
	}

	return app
}

func logFileServerError(cfg *Config, err error) {
	if err == nil || cfg == nil || cfg.Log == nil {
		return
	}
	cfg.Log.Error(context.Background(), "fileserver", "err", err)
}
