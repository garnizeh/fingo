package homebus

import (
	"context"
	cryptorand "crypto/rand"
	"fmt"
	"math/big"

	"github.com/garnizeh/fingo/business/types/home"
	"github.com/google/uuid"
)

// TestGenerateNewHomes is a helper method for testing.
func TestGenerateNewHomes(n int, userID uuid.UUID) []NewHome {
	newHmes := make([]NewHome, n)

	idx := randomInt(10000)
	for i := range newHmes {
		idx++

		t := home.Single
		if v := (idx + i) % 2; v == 0 {
			t = home.Condo
		}

		address := Address{
			Address1: fmt.Sprintf("Address%d", idx),
			Address2: fmt.Sprintf("Address%d", idx),
			ZipCode:  fmt.Sprintf("%05d", idx),
			City:     fmt.Sprintf("City%d", idx),
			State:    fmt.Sprintf("State%d", idx),
			Country:  fmt.Sprintf("Country%d", idx),
		}
		nh := NewHome{
			Type:    t,
			Address: &address,
			UserID:  userID,
		}

		newHmes[i] = nh
	}

	return newHmes
}

func randomInt(limit int) int {
	if limit <= 0 {
		return 0
	}
	n, err := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(limit)))
	if err != nil {
		panic(err)
	}
	return int(n.Int64())
}

// TestGenerateSeedHomes is a helper method for testing.
func TestGenerateSeedHomes(ctx context.Context, n int, api ExtBusiness, userID uuid.UUID) ([]Home, error) {
	newHmes := TestGenerateNewHomes(n, userID)

	hmes := make([]Home, len(newHmes))
	for i := range newHmes {
		nh := &newHmes[i]
		hme, err := api.Create(ctx, nh)
		if err != nil {
			return nil, fmt.Errorf("seeding home: idx: %d : %w", i, err)
		}

		hmes[i] = hme
	}

	return hmes, nil
}

// ParseAddress is a helper function to create an address value.
func ParseAddress(address1, address2, zipCode, city, state, country string) Address {
	return Address{
		Address1: address1,
		Address2: address2,
		ZipCode:  zipCode,
		City:     city,
		State:    state,
		Country:  country,
	}
}
