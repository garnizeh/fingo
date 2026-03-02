package apitest

import (
	"net/http/httptest"
	"testing"

	authbuild "github.com/garnizeh/fingo/api/services/auth/build"
	salesbuild "github.com/garnizeh/fingo/api/services/fingo/build"
	"github.com/garnizeh/fingo/app/sdk/auth"
	"github.com/garnizeh/fingo/app/sdk/authclient/http"
	"github.com/garnizeh/fingo/app/sdk/mux"
	"github.com/garnizeh/fingo/business/sdk/dbtest"
)

// New initialized the system to run a test.
func New(t *testing.T, testName string) *Test {
	db := dbtest.New(t, testName)

	// -------------------------------------------------------------------------

	auth := auth.New(auth.Config{
		Log:       db.Log,
		UserBus:   db.BusDomain.User,
		KeyLookup: &KeyStore{},
	})

	// -------------------------------------------------------------------------

	handler := mux.WebAPI(&mux.Config{
		Log: db.Log,
		DB:  db.DB,
		BusConfig: mux.BusConfig{
			UserBus: db.BusDomain.User,
		},
		AuthConfig: mux.AuthConfig{
			Auth: auth,
		},
	}, authbuild.Routes())
	server := httptest.NewServer(handler)

	authClient, err := http.New(db.Log, server.URL)
	if err != nil {
		t.Fatal("could not create authentication client")
	}

	// -------------------------------------------------------------------------

	mux := mux.WebAPI(&mux.Config{
		Log: db.Log,
		DB:  db.DB,
		BusConfig: mux.BusConfig{
			AuditBus:    db.BusDomain.Audit,
			UserBus:     db.BusDomain.User,
			ProductBus:  db.BusDomain.Product,
			HomeBus:     db.BusDomain.Home,
			VProductBus: db.BusDomain.VProduct,
		},
		FinGoConfig: mux.FinGoConfig{
			AuthClient: authClient,
		},
	}, salesbuild.Routes())

	return &Test{
		DB:   db,
		Auth: auth,
		mux:  mux,
	}
}
