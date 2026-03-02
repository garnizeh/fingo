package main

import (
	"context"
	"embed"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/garnizeh/fingo/api/services/fingo/build"
	"github.com/garnizeh/fingo/app/sdk/authclient"
	"github.com/garnizeh/fingo/app/sdk/authclient/grpc"
	http2 "github.com/garnizeh/fingo/app/sdk/authclient/http"
	"github.com/garnizeh/fingo/app/sdk/debug"
	"github.com/garnizeh/fingo/app/sdk/mux"
	"github.com/garnizeh/fingo/business/domain/auditbus"
	"github.com/garnizeh/fingo/business/domain/auditbus/extensions/auditotel"
	"github.com/garnizeh/fingo/business/domain/auditbus/stores/auditdb"
	"github.com/garnizeh/fingo/business/domain/creditcardbus"
	"github.com/garnizeh/fingo/business/domain/creditcardbus/extensions/creditcardaudit"
	"github.com/garnizeh/fingo/business/domain/creditcardbus/extensions/creditcardotel"
	"github.com/garnizeh/fingo/business/domain/creditcardbus/stores/creditcarddb"
	"github.com/garnizeh/fingo/business/domain/homebus"
	"github.com/garnizeh/fingo/business/domain/homebus/extensions/homeotel"
	"github.com/garnizeh/fingo/business/domain/homebus/stores/homedb"
	"github.com/garnizeh/fingo/business/domain/productbus"
	"github.com/garnizeh/fingo/business/domain/productbus/extensions/productotel"
	"github.com/garnizeh/fingo/business/domain/productbus/stores/productdb"
	"github.com/garnizeh/fingo/business/domain/userbus"
	"github.com/garnizeh/fingo/business/domain/userbus/extensions/useraudit"
	"github.com/garnizeh/fingo/business/domain/userbus/extensions/userotel"
	"github.com/garnizeh/fingo/business/domain/userbus/stores/usercache"
	"github.com/garnizeh/fingo/business/domain/userbus/stores/userdb"
	"github.com/garnizeh/fingo/business/domain/vproductbus"
	"github.com/garnizeh/fingo/business/domain/vproductbus/extensions/vproductotel"
	"github.com/garnizeh/fingo/business/domain/vproductbus/stores/vproductdb"
	"github.com/garnizeh/fingo/business/sdk/delegate"
	"github.com/garnizeh/fingo/business/sdk/sqldb"
	"github.com/garnizeh/fingo/foundation/logger"
	"github.com/garnizeh/fingo/foundation/otel"
)

/*
	Need to figure out timeouts for http service.
*/

//go:embed static
var static embed.FS

var tag = "develop"

func main() {
	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT *******")
		},
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "FINGO", otel.GetTraceID, events)

	// -------------------------------------------------------------------------

	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	// -------------------------------------------------------------------------
	// GOMAXPROCS

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// Configuration

	//nolint:govet // Configuration layout is intentionally organized for conf parsing clarity.
	cfg := struct {
		conf.Version
		Web struct {
			APIHost            string        `conf:"default:0.0.0.0:3000"`
			DebugHost          string        `conf:"default:0.0.0.0:3010"`
			CORSAllowedOrigins []string      `conf:"default:*"`
			ReadTimeout        time.Duration `conf:"default:5s"`
			WriteTimeout       time.Duration `conf:"default:10s"`
			IdleTimeout        time.Duration `conf:"default:120s"`
			ShutdownTimeout    time.Duration `conf:"default:20s"`
		}
		Auth struct {
			Host      string `conf:"default:http://auth-service:6000"`
			Mechanism string `conf:"default:http"`
			GRPC      struct {
				Host string `conf:"default:auth-service:6001"`
			}
		}
		DB struct {
			User         string `conf:"default:postgres"`
			Password     string `conf:"default:postgres,mask"`
			Host         string `conf:"default:database-service"`
			Name         string `conf:"default:postgres"`
			MaxIdleConns int    `conf:"default:0"`
			MaxOpenConns int    `conf:"default:0"`
			DisableTLS   bool   `conf:"default:true"`
		}
		Tempo struct {
			Host        string  `conf:"default:tempo:4317"`
			ServiceName string  `conf:"default:fingo"`
			Probability float64 `conf:"default:0.05"`
			// Shouldn't use a high Probability value in non-developer systems.
			// 0.05 should be enough for most systems. Some might want to have
			// this even lower.
		}
	}{
		Version: conf.Version{
			Build: tag,
			Desc:  "FinGo",
		},
	}

	const prefix = "FINGO"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	// -------------------------------------------------------------------------
	// App Starting

	log.Info(ctx, "starting service", "version", cfg.Build)
	defer log.Info(ctx, "shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Info(ctx, "startup", "config", out)

	log.BuildInfo(ctx)

	expvar.NewString("build").Set(cfg.Build)

	// -------------------------------------------------------------------------
	// Database Support

	log.Info(ctx, "startup", "status", "initializing database support", "hostport", cfg.DB.Host)

	dbcfg := sqldb.Config{
		User:         cfg.DB.User,
		Password:     cfg.DB.Password,
		Host:         cfg.DB.Host,
		Name:         cfg.DB.Name,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxOpenConns: cfg.DB.MaxOpenConns,
		DisableTLS:   cfg.DB.DisableTLS,
	}
	db, err := sqldb.Open(&dbcfg)
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}

	defer func() {
		if errClose := db.Close(); errClose != nil {
			log.Error(ctx, "db close", "err", errClose)
		}
	}()

	// -------------------------------------------------------------------------
	// Create Business Packages

	delegate := delegate.New(log)

	auditOtelExt := auditotel.NewExtension()
	auditStorage := auditdb.NewStore(log, db)
	auditBus := auditbus.NewBusiness(log, auditStorage, auditOtelExt)

	userOtelExt := userotel.NewExtension()
	userAuditExt := useraudit.NewExtension(auditBus)
	userStorage := usercache.NewStore(log, userdb.NewStore(log, db), time.Minute)
	userBus := userbus.NewBusiness(log, delegate, userStorage, userOtelExt, userAuditExt)

	creditCardOtelExt := creditcardotel.NewExtension()
	creditCardAuditExt := creditcardaudit.NewExtension(auditBus)
	creditCardStorage := creditcarddb.NewStore(log, db)
	creditCardBus := creditcardbus.NewBusiness(log, userBus, delegate, creditCardStorage, creditCardOtelExt, creditCardAuditExt)

	productOtelExt := productotel.NewExtension()
	productStorage := productdb.NewStore(log, db)
	productBus := productbus.NewBusiness(log, userBus, delegate, productStorage, productOtelExt)

	homeOtelExt := homeotel.NewExtension()
	homeStorage := homedb.NewStore(log, db)
	homeBus := homebus.NewBusiness(log, userBus, delegate, homeStorage, homeOtelExt)

	vproductOtelExt := vproductotel.NewExtension()
	vproductStorage := vproductdb.NewStore(log, db)
	vproductBus := vproductbus.NewBusiness(vproductStorage, vproductOtelExt)

	// -------------------------------------------------------------------------
	// Initialize authentication support

	log.Info(ctx, "startup", "status", "initializing authentication support")

	var authClient authclient.Authenticator
	switch cfg.Auth.Mechanism {
	case "grpc":
		authClient, err = grpc.New(log, cfg.Auth.GRPC.Host)
	default:
		authClient, err = http2.New(log, cfg.Auth.Host)
	}

	if err != nil {
		log.Error(ctx, "failed to initialize authentication client", "error", err)
		return fmt.Errorf("failed to initialize authentication client: %w", err)
	}

	defer func() {
		if errClose := authClient.Close(); errClose != nil {
			log.Error(ctx, "authclient close", "err", errClose)
		}
	}()

	// -------------------------------------------------------------------------
	// Start Tracing Support

	log.Info(ctx, "startup", "status", "initializing tracing support")

	traceProvider, teardown, err := otel.InitTracing(otel.Config{
		ServiceName: cfg.Tempo.ServiceName,
		Host:        cfg.Tempo.Host,
		ExcludedRoutes: map[string]struct{}{
			"/v1/liveness":  {},
			"/v1/readiness": {},
		},
		Probability: cfg.Tempo.Probability,
	})
	if err != nil {
		return fmt.Errorf("starting tracing: %w", err)
	}

	defer teardown(context.Background())

	tracer := traceProvider.Tracer(cfg.Tempo.ServiceName)

	// -------------------------------------------------------------------------
	// Start Debug Service

	go func() {
		log.Info(ctx, "startup", "status", "debug v1 router started", "host", cfg.Web.DebugHost)

		srv := http.Server{
			Addr:              cfg.Web.DebugHost,
			Handler:           debug.Mux(),
			ReadHeaderTimeout: 5 * time.Second,
		}

		if err := srv.ListenAndServe(); err != nil {
			log.Error(ctx, "shutdown", "status", "debug v1 router closed", "host", cfg.Web.DebugHost, "msg", err)
		}
	}()

	// -------------------------------------------------------------------------
	// Start API Service

	log.Info(ctx, "startup", "status", "initializing V1 API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	cfgMux := mux.Config{
		Build:  tag,
		Log:    log,
		DB:     db,
		Tracer: tracer,
		BusConfig: mux.BusConfig{
			AuditBus:      auditBus,
			UserBus:       userBus,
			CreditCardBus: creditCardBus,
			ProductBus:    productBus,
			HomeBus:       homeBus,
			VProductBus:   vproductBus,
		},
		FinGoConfig: mux.FinGoConfig{
			AuthClient: authClient,
		},
	}

	webAPI := mux.WebAPI(&cfgMux, build.Routes(), mux.WithCORS(cfg.Web.CORSAllowedOrigins), mux.WithFileServer(false, static, "static", "/"))

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      webAPI,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", api.Addr)

		serverErrors <- api.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		ctxTimeout, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctxTimeout); err != nil {
			if errClose := api.Close(); errClose != nil {
				log.Error(ctx, "api close", "err", errClose)
			}
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
