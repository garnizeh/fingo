// Package creditcarddb contains credit card related CRUD functionality.
package creditcarddb

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/garnizeh/fingo/business/domain/creditcardbus"
	"github.com/garnizeh/fingo/business/sdk/order"
	"github.com/garnizeh/fingo/business/sdk/page"
	"github.com/garnizeh/fingo/business/sdk/sqldb"
	"github.com/garnizeh/fingo/foundation/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Store manages the set of APIs for credit card database access.
type Store struct {
	log *logger.Logger
	db  sqlx.ExtContext
}

// NewStore constructs the api for data access.
func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

// NewWithTx constructs a new Store value replacing the sqlx DB
// value with a sqlx DB value that is currently inside a transaction.
func (s *Store) NewWithTx(tx sqldb.CommitRollbacker) (creditcardbus.Storer, error) {
	ec, err := sqldb.GetExtContext(tx)
	if err != nil {
		return nil, err
	}

	return &Store{
		log: s.log,
		db:  ec,
	}, nil
}

// Create adds a CreditCard to the sqldb.
func (s *Store) Create(ctx context.Context, cc creditcardbus.CreditCard) error {
	const q = `
	INSERT INTO credit_cards
		(credit_card_id, user_id, name, card_limit, closing_day, due_day, last_four_digits, enabled, date_created, date_updated)
	VALUES
		(:credit_card_id, :user_id, :name, :card_limit, :closing_day, :due_day, :last_four_digits, :enabled, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBCreditCard(cc)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update modifies data about a CreditCard.
func (s *Store) Update(ctx context.Context, cc creditcardbus.CreditCard) error {
	const q = `
	UPDATE
		credit_cards
	SET
		"name" = :name,
		"card_limit" = :card_limit,
		"closing_day" = :closing_day,
		"due_day" = :due_day,
		"enabled" = :enabled,
		"date_updated" = :date_updated
	WHERE
		credit_card_id = :credit_card_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBCreditCard(cc)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Delete removes the credit card identified by a given ID.
func (s *Store) Delete(ctx context.Context, cc creditcardbus.CreditCard) error {
	const q = `
	DELETE FROM
		credit_cards
	WHERE
		credit_card_id = :credit_card_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBCreditCard(cc)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Query gets all CreditCards from the database.
func (s *Store) Query(ctx context.Context, filter creditcardbus.QueryFilter, orderBy order.By, page page.Page) ([]creditcardbus.CreditCard, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		credit_card_id, user_id, name, card_limit, closing_day, due_day, last_four_digits, enabled, date_created, date_updated
	FROM
		credit_cards`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbCCs []dbCreditCard
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbCCs); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusCreditCards(dbCCs)
}

// Count returns the total number of credit cards in the database.
func (s *Store) Count(ctx context.Context, filter creditcardbus.QueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		credit_cards`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("namedquerystruct: %w", err)
	}

	return count.Count, nil
}

// QueryByID finds the credit card identified by a given ID.
func (s *Store) QueryByID(ctx context.Context, ccID uuid.UUID) (creditcardbus.CreditCard, error) {
	data := struct {
		ID string `db:"credit_card_id"`
	}{
		ID: ccID.String(),
	}

	const q = `
	SELECT
		credit_card_id, user_id, name, card_limit, closing_day, due_day, last_four_digits, enabled, date_created, date_updated
	FROM
		credit_cards
	WHERE
		credit_card_id = :credit_card_id`

	var dbCC dbCreditCard
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbCC); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) || errors.Is(err, sql.ErrNoRows) {
			return creditcardbus.CreditCard{}, fmt.Errorf("namedquerystruct: %w", creditcardbus.ErrNotFound)
		}
		return creditcardbus.CreditCard{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	bus, err := toBusCreditCard(dbCC)
	if err != nil {
		return creditcardbus.CreditCard{}, fmt.Errorf("toBusCreditCard: %w", err)
	}

	return bus, nil
}
