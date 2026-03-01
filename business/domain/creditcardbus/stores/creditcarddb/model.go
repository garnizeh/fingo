package creditcarddb

import (
	"fmt"
	"time"

	"github.com/garnizeh/fingo/business/domain/creditcardbus"
	"github.com/garnizeh/fingo/business/types/money"
	"github.com/garnizeh/fingo/business/types/name"
	"github.com/google/uuid"
)

type dbCreditCard struct {
	ID             uuid.UUID `db:"credit_card_id"`
	UserID         uuid.UUID `db:"user_id"`
	Name           string    `db:"name"`
	Limit          float64   `db:"card_limit"`
	ClosingDay     int       `db:"closing_day"`
	DueDay         int       `db:"due_day"`
	LastFourDigits string    `db:"last_four_digits"`
	Enabled        bool      `db:"enabled"`
	DateCreated    time.Time `db:"date_created"`
	DateUpdated    time.Time `db:"date_updated"`
}

func toDBCreditCard(bus creditcardbus.CreditCard) dbCreditCard {
	return dbCreditCard{
		ID:             bus.ID,
		UserID:         bus.UserID,
		Name:           bus.Name.String(),
		Limit:          bus.Limit.Value(),
		ClosingDay:     bus.ClosingDay,
		DueDay:         bus.DueDay,
		LastFourDigits: bus.LastFourDigits,
		Enabled:        bus.Enabled,
		DateCreated:    bus.DateCreated.UTC(),
		DateUpdated:    bus.DateUpdated.UTC(),
	}
}

func toBusCreditCard(db dbCreditCard) (creditcardbus.CreditCard, error) {
	name, err := name.Parse(db.Name)
	if err != nil {
		return creditcardbus.CreditCard{}, fmt.Errorf("parse name: %w", err)
	}

	limit, err := money.Parse(db.Limit)
	if err != nil {
		return creditcardbus.CreditCard{}, fmt.Errorf("parse limit: %w", err)
	}

	bus := creditcardbus.CreditCard{
		ID:             db.ID,
		UserID:         db.UserID,
		Name:           name,
		Limit:          limit,
		ClosingDay:     db.ClosingDay,
		DueDay:         db.DueDay,
		LastFourDigits: db.LastFourDigits,
		Enabled:        db.Enabled,
		DateCreated:    db.DateCreated.In(time.Local),
		DateUpdated:    db.DateUpdated.In(time.Local),
	}

	return bus, nil
}

func toBusCreditCards(dbs []dbCreditCard) ([]creditcardbus.CreditCard, error) {
	bus := make([]creditcardbus.CreditCard, len(dbs))

	for i, db := range dbs {
		var err error
		bus[i], err = toBusCreditCard(db)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}
