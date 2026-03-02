package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/garnizeh/fingo/business/sdk/migrate"
	"github.com/garnizeh/fingo/business/sdk/sqldb"
)

// Seed loads test data into the database.
func Seed(cfg *sqldb.Config) (err error) {
	db, err := sqldb.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			err = errors.Join(err, fmt.Errorf("closing database: %w", cerr))
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := migrate.Seed(ctx, db); err != nil {
		return fmt.Errorf("seed database: %w", err)
	}

	fmt.Println("seed data complete")
	return nil
}
