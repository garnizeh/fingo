package unittest

import (
	"context"

	"github.com/garnizeh/fingo/business/domain/auditbus"
	"github.com/garnizeh/fingo/business/domain/homebus"
	"github.com/garnizeh/fingo/business/domain/productbus"
	"github.com/garnizeh/fingo/business/domain/userbus"
)

// User represents an app user specified for the test.
type User struct {
	Products []productbus.Product
	Homes    []homebus.Home
	Audits   []auditbus.Audit
	userbus.User
}

// SeedData represents data that was seeded for the test.
type SeedData struct {
	Users  []User
	Admins []User
}

// Table represent fields needed for running an unit test.
type Table struct {
	ExpResp any
	ExcFunc func(ctx context.Context) any
	CmpFunc func(got any, exp any) string
	Name    string
}
