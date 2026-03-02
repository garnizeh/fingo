package apitest

import (
	"github.com/garnizeh/fingo/business/domain/auditbus"
	"github.com/garnizeh/fingo/business/domain/homebus"
	"github.com/garnizeh/fingo/business/domain/productbus"
	"github.com/garnizeh/fingo/business/domain/userbus"
)

// User extends the dbtest user for api test support.
type User struct {
	Token    string
	Products []productbus.Product
	Homes    []homebus.Home
	Audits   []auditbus.Audit
	userbus.User
}

// SeedData represents users for api tests.
type SeedData struct {
	Users  []User
	Admins []User
}

// Table represent fields needed for running an api test.
type Table struct {
	Input      any
	GotResp    any
	ExpResp    any
	CmpFunc    func(got any, exp any) string
	Name       string
	URL        string
	Token      string
	Method     string
	StatusCode int
}
