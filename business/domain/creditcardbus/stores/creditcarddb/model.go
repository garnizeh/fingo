package creditcarddb

import (
	"fmt"
	"time"

	"github.com/garnizeh/fingo/business/domain/creditcardbus"
	"github.com/garnizeh/fingo/business/types/money"
	"github.com/garnizeh/fingo/business/types/name"
	"github.com/google/uuid"
)

type dbCreditCardTimestamps struct {
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

type dbCreditCardIdentity struct {
	Name           string `db:"name"`
	LastFourDigits string `db:"last_four_digits"`
}

type dbCreditCard struct {
	dbCreditCardTimestamps
	dbCreditCardIdentity
	ID         uuid.UUID `db:"credit_card_id"`
	UserID     uuid.UUID `db:"user_id"`
	Limit      float64   `db:"card_limit"`
	ClosingDay int       `db:"closing_day"`
	DueDay     int       `db:"due_day"`
	Enabled    bool      `db:"enabled"`
}

func toDBCreditCard(bus *creditcardbus.CreditCard) dbCreditCard {
	if bus == nil {
		return dbCreditCard{}
	}
	return dbCreditCard{
		dbCreditCardTimestamps: dbCreditCardTimestamps{
			DateCreated: bus.DateCreated.UTC(),
			DateUpdated: bus.DateUpdated.UTC(),
		},
		dbCreditCardIdentity: dbCreditCardIdentity{
			Name:           bus.Name.String(),
			LastFourDigits: bus.LastFourDigits,
		},
		ID:         bus.ID,
		UserID:     bus.UserID,
		Limit:      bus.Limit.Value(),
		ClosingDay: bus.ClosingDay,
		DueDay:     bus.DueDay,
		Enabled:    bus.Enabled,
	}
}

func toBusCreditCard(db *dbCreditCard) (creditcardbus.CreditCard, error) {
	name, err := name.Parse(db.Name)
	if err != nil {
		return creditcardbus.CreditCard{}, fmt.Errorf("parse name: %w", err)
	}

	limit, err := money.Parse(db.Limit)
	if err != nil {
		return creditcardbus.CreditCard{}, fmt.Errorf("parse limit: %w", err)
	}

	bus := creditcardbus.CreditCard{
		CreditCardTimestamps: creditcardbus.CreditCardTimestamps{
			DateCreated: db.DateCreated.In(time.Local),
			DateUpdated: db.DateUpdated.In(time.Local),
		},
		CreditCardIdentity: creditcardbus.CreditCardIdentity{
			Name:           name,
			LastFourDigits: db.LastFourDigits,
		},
		ID:         db.ID,
		UserID:     db.UserID,
		Limit:      limit,
		ClosingDay: db.ClosingDay,
		DueDay:     db.DueDay,
		Enabled:    db.Enabled,
	}

	return bus, nil
}

func toBusCreditCards(dbs []dbCreditCard) ([]creditcardbus.CreditCard, error) {
	bus := make([]creditcardbus.CreditCard, len(dbs))

	for i := range dbs {
		var err error
		bus[i], err = toBusCreditCard(&dbs[i])
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}
