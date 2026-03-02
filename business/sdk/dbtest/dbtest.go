// Package dbtest contains supporting code for running tests that hit the DB.
package dbtest

import (
	"bytes"
	"context"
	cryptorand "crypto/rand"
	"math/big"
	"net/url"
	"testing"
	"time"

	"github.com/garnizeh/fingo/business/sdk/migrate"
	"github.com/garnizeh/fingo/business/sdk/sqldb"
	"github.com/garnizeh/fingo/foundation/logger"
	"github.com/garnizeh/fingo/foundation/otel"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// Database owns state for running and shutting down tests.
type Database struct {
	DB        *sqlx.DB
	Log       *logger.Logger
	BusDomain BusDomain
}

// New creates a new test database inside the database that was started
// to handle testing. The database is migrated to the current version and
// a connection pool is provided with business domain packages.
func New(t *testing.T, testName string) *Database {
	t.Helper()

	ctx := t.Context()

	const (
		image    = "postgres:18.2-alpine"
		database = "fingo_test"
		username = "postgres"
		password = "postgres"
	)

	pgContainer, err := postgres.Run(
		ctx,
		image,
		postgres.WithDatabase(database),
		postgres.WithUsername(username),
		postgres.WithPassword(password),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("connection string: %v", err)
	}

	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("parse dsn: %v", err)
	}

	response, err := pgContainer.Inspect(ctx)
	if err != nil {
		t.Fatalf("inspect container: %v", err)
	}

	t.Logf("Name    : %s\n", response.Name)
	t.Logf("HostPort: %s\n", u.Host)

	cfgdbM := sqldb.Config{
		User:       username,
		Password:   password,
		Host:       u.Host,
		Name:       database,
		DisableTLS: true,
	}
	dbM, err := sqldb.Open(&cfgdbM)
	if err != nil {
		t.Fatalf("Opening database connection: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if serr := sqldb.StatusCheck(ctx, dbM); serr != nil {
		t.Fatalf("status check database: %v", serr)
	}

	// -------------------------------------------------------------------------

	const letterBytes = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, 4)
	for i := range b {
		b[i] = letterBytes[randomInt(len(letterBytes))]
	}
	dbName := string(b)

	t.Logf("Create Database: %s\n", dbName)
	if _, eerr := dbM.ExecContext(context.Background(), "CREATE DATABASE "+dbName); eerr != nil {
		t.Fatalf("creating database %s: %v", dbName, eerr)
	}

	// -------------------------------------------------------------------------

	cfgdb := sqldb.Config{
		User:       username,
		Password:   password,
		Host:       u.Host,
		Name:       dbName,
		DisableTLS: true,
	}
	db, err := sqldb.Open(&cfgdb)
	if err != nil {
		t.Fatalf("Opening database connection: %v", err)
	}

	t.Logf("Migrate Database: %s\n", dbName)
	if err := migrate.Migrate(ctx, db); err != nil {
		t.Fatalf("Migrating error: %s", err)
	}

	// -------------------------------------------------------------------------

	var buf bytes.Buffer
	log := logger.New(&buf, logger.LevelInfo, "TEST", func(context.Context) string { return otel.GetTraceID(ctx) })

	// -------------------------------------------------------------------------

	t.Cleanup(func() {
		t.Helper()

		t.Logf("Drop Database: %s\n", dbName)
		if _, err := dbM.ExecContext(context.Background(), "DROP DATABASE "+dbName); err != nil {
			t.Fatalf("dropping database %s: %v", dbName, err)
		}

		if err := db.Close(); err != nil {
			t.Fatalf("closing db: %v", err)
		}

		if err := dbM.Close(); err != nil {
			t.Fatalf("closing db manager: %v", err)
		}

		t.Logf("******************** LOGS (%s) ********************\n\n", testName)
		t.Log(buf.String())
		t.Logf("******************** LOGS (%s) ********************\n", testName)
	})

	return &Database{
		DB:        db,
		Log:       log,
		BusDomain: newBusDomains(log, db),
	}
}

func randomInt(maxValue int) int {
	if maxValue <= 0 {
		return 0
	}
	n, err := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(maxValue)))
	if err != nil {
		panic(err)
	}
	return int(n.Int64())
}
