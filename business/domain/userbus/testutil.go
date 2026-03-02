package userbus

import (
	"context"
	cryptorand "crypto/rand"
	"fmt"
	"math/big"
	"net/mail"

	"github.com/garnizeh/fingo/business/types/name"
	"github.com/garnizeh/fingo/business/types/password"
	"github.com/garnizeh/fingo/business/types/role"
	"github.com/google/uuid"
)

// TestNewUsers is a helper method for testing.
func TestNewUsers(n int, rle role.Role) []NewUser {
	newUsrs := make([]NewUser, n)

	idx := randomInt(10000)
	for i := range n {
		idx++

		nu := NewUser{
			Name:       name.MustParse(fmt.Sprintf("Name%d", idx)),
			Email:      mail.Address{Address: fmt.Sprintf("Email%d@gmail.com", idx)},
			Roles:      []role.Role{rle},
			Department: name.MustParseNull(fmt.Sprintf("Department%d", idx)),
			Password:   password.MustParse(fmt.Sprintf("Password%d", idx)),
		}

		newUsrs[i] = nu
	}

	return newUsrs
}

func randomInt(n int) int {
	if n <= 0 {
		return 0
	}
	res, err := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(n)))
	if err != nil {
		panic(err)
	}
	return int(res.Int64())
}

// TestSeedUsers is a helper method for testing.
func TestSeedUsers(ctx context.Context, n int, role role.Role, api ExtBusiness) ([]User, error) {
	newUsrs := TestNewUsers(n, role)

	usrs := make([]User, len(newUsrs))
	for i, nu := range newUsrs {
		usr, err := api.Create(ctx, uuid.UUID{}, &nu)
		if err != nil {
			return nil, fmt.Errorf("seeding user: idx: %d : %w", i, err)
		}

		usrs[i] = usr
	}

	return usrs, nil
}
