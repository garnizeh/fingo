package homedb

import (
	"fmt"
	"time"

	"github.com/garnizeh/fingo/business/domain/homebus"
	"github.com/garnizeh/fingo/business/types/home"
	"github.com/google/uuid"
)

type homeDB struct {
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
	Type        string    `db:"type"`
	Address1    string    `db:"address_1"`
	Address2    string    `db:"address_2"`
	ZipCode     string    `db:"zip_code"`
	City        string    `db:"city"`
	Country     string    `db:"country"`
	State       string    `db:"state"`
	ID          uuid.UUID `db:"home_id"`
	UserID      uuid.UUID `db:"user_id"`
}

func toDBHome(bus *homebus.Home) homeDB {
	db := homeDB{
		ID:          bus.ID,
		UserID:      bus.UserID,
		Type:        bus.Type.String(),
		DateCreated: bus.DateCreated.UTC(),
		DateUpdated: bus.DateUpdated.UTC(),
	}

	if bus.Address != nil {
		db.Address1 = bus.Address.Address1
		db.Address2 = bus.Address.Address2
		db.ZipCode = bus.Address.ZipCode
		db.City = bus.Address.City
		db.Country = bus.Address.Country
		db.State = bus.Address.State
	}

	return db
}

func toBusHome(db *homeDB) (homebus.Home, error) {
	typ, err := home.Parse(db.Type)
	if err != nil {
		return homebus.Home{}, fmt.Errorf("parse type: %w", err)
	}

	return homebus.Home{
		ID:     db.ID,
		UserID: db.UserID,
		Type:   typ,
		Address: &homebus.Address{
			Address1: db.Address1,
			Address2: db.Address2,
			ZipCode:  db.ZipCode,
			City:     db.City,
			Country:  db.Country,
			State:    db.State,
		},
		DateCreated: db.DateCreated.In(time.Local),
		DateUpdated: db.DateUpdated.In(time.Local),
	}, nil
}

func toBusHomes(dbs []homeDB) ([]homebus.Home, error) {
	bus := make([]homebus.Home, len(dbs))
	for i := range dbs {
		item, err := toBusHome(&dbs[i])
		if err != nil {
			return nil, fmt.Errorf("parse type: %w", err)
		}
		bus[i] = item
	}

	return bus, nil
}
