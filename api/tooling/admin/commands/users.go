package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/garnizeh/fingo/business/domain/userbus"
	"github.com/garnizeh/fingo/business/domain/userbus/stores/userdb"
	"github.com/garnizeh/fingo/business/sdk/page"
	"github.com/garnizeh/fingo/business/sdk/sqldb"
	"github.com/garnizeh/fingo/foundation/logger"
)

// Users retrieves all users from the database.
func Users(log *logger.Logger, cfg *sqldb.Config, pageNumber, rowsPerPage string) (err error) {
	db, err := sqldb.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			err = errors.Join(err, fmt.Errorf("closing database: %w", cerr))
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userBus := userbus.NewBusiness(log, nil, userdb.NewStore(log, db))

	page, err := page.Parse(pageNumber, rowsPerPage)
	if err != nil {
		return fmt.Errorf("parsing page information: %w", err)
	}

	users, err := userBus.Query(ctx, userbus.QueryFilter{}, userbus.DefaultOrderBy, page)
	if err != nil {
		return fmt.Errorf("retrieve users: %w", err)
	}

	return json.NewEncoder(os.Stdout).Encode(users)
}
