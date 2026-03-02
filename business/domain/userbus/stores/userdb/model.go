package userdb

import (
	"database/sql"
	"fmt"
	"net/mail"
	"time"

	"github.com/garnizeh/fingo/business/domain/userbus"
	"github.com/garnizeh/fingo/business/sdk/sqldb/dbarray"
	"github.com/garnizeh/fingo/business/types/name"
	"github.com/garnizeh/fingo/business/types/role"
	"github.com/google/uuid"
)

type userDB struct {
	DateCreated  time.Time      `db:"date_created"`
	DateUpdated  time.Time      `db:"date_updated"`
	Name         string         `db:"name"`
	Email        string         `db:"email"`
	Roles        dbarray.String `db:"roles"`
	PasswordHash []byte         `db:"password_hash"`
	Department   sql.NullString `db:"department"`
	ID           uuid.UUID      `db:"user_id"`
	Enabled      bool           `db:"enabled"`
}

func toDBUser(bus *userbus.User) userDB {
	return userDB{
		ID:           bus.ID,
		Name:         bus.Name.String(),
		Email:        bus.Email.Address,
		Roles:        role.ParseToString(bus.Roles),
		PasswordHash: bus.PasswordHash,
		Department:   name.ToSQLNullString(bus.Department),
		Enabled:      bus.Enabled,
		DateCreated:  bus.DateCreated.UTC(),
		DateUpdated:  bus.DateUpdated.UTC(),
	}
}

func toBusUser(db *userDB) (userbus.User, error) {
	addr := mail.Address{
		Address: db.Email,
	}

	roles, err := role.ParseMany(db.Roles)
	if err != nil {
		return userbus.User{}, fmt.Errorf("parse: %w", err)
	}

	nme, err := name.Parse(db.Name)
	if err != nil {
		return userbus.User{}, fmt.Errorf("parse name: %w", err)
	}

	department, err := name.ParseNull(db.Department.String)
	if err != nil {
		return userbus.User{}, fmt.Errorf("parse department: %w", err)
	}

	bus := userbus.User{
		ID:           db.ID,
		Name:         nme,
		Email:        addr,
		Roles:        roles,
		PasswordHash: db.PasswordHash,
		Enabled:      db.Enabled,
		Department:   department,
		DateCreated:  db.DateCreated.In(time.Local),
		DateUpdated:  db.DateUpdated.In(time.Local),
	}

	return bus, nil
}

func toBusUsers(dbs []userDB) ([]userbus.User, error) {
	bus := make([]userbus.User, len(dbs))

	for i := range dbs {
		var err error
		bus[i], err = toBusUser(&dbs[i])
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}
