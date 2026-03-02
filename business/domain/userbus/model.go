package userbus

import (
	"net/mail"
	"time"

	"github.com/garnizeh/fingo/business/types/name"
	"github.com/garnizeh/fingo/business/types/password"
	"github.com/garnizeh/fingo/business/types/role"
	"github.com/google/uuid"
)

// User represents information about an individual user.
type User struct {
	DateCreated  time.Time
	DateUpdated  time.Time
	Email        mail.Address
	Name         name.Name
	Roles        []role.Role
	PasswordHash []byte
	Department   name.Null
	ID           uuid.UUID
	Enabled      bool
}

// NewUser contains information needed to create a new user.
type NewUser struct {
	Email      mail.Address
	Name       name.Name
	Password   password.Password
	Roles      []role.Role
	Department name.Null
}

// UpdateUser contains information needed to update a user.
type UpdateUser struct {
	Name       *name.Name
	Email      *mail.Address
	Department *name.Null
	Password   *password.Password
	Enabled    *bool
	Roles      []role.Role
}
