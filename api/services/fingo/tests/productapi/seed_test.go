package product_test

import (
	"context"
	"fmt"

	"github.com/garnizeh/fingo/app/sdk/apitest"
	"github.com/garnizeh/fingo/app/sdk/auth"
	"github.com/garnizeh/fingo/business/domain/productbus"
	"github.com/garnizeh/fingo/business/domain/userbus"
	"github.com/garnizeh/fingo/business/sdk/dbtest"
	"github.com/garnizeh/fingo/business/types/role"
)

func insertSeedData(db *dbtest.Database, ath *auth.Auth) (apitest.SeedData, error) {
	ctx := context.Background()
	busDomain := db.BusDomain

	usrs, err := userbus.TestSeedUsers(ctx, 1, role.User, busDomain.User)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	prds, err := productbus.TestGenerateSeedProducts(ctx, 2, busDomain.Product, usrs[0].ID)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding products : %w", err)
	}

	tu1 := apitest.User{
		User:     usrs[0],
		Products: prds,

		Token: apitest.Token(db.BusDomain.User, ath, usrs[0].Email.Address),
	}

	// -------------------------------------------------------------------------

	usrs, err = userbus.TestSeedUsers(ctx, 1, role.Admin, busDomain.User)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	prds, err = productbus.TestGenerateSeedProducts(ctx, 2, busDomain.Product, usrs[0].ID)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding products : %w", err)
	}

	tu2 := apitest.User{
		User:     usrs[0],
		Products: prds,
		Token:    apitest.Token(db.BusDomain.User, ath, usrs[0].Email.Address),
	}

	// -------------------------------------------------------------------------

	sd := apitest.SeedData{
		Admins: []apitest.User{tu2},
		Users:  []apitest.User{tu1},
	}

	return sd, nil
}
