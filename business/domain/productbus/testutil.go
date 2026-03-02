package productbus

import (
	"context"
	cryptorand "crypto/rand"
	"fmt"
	"math/big"

	"github.com/garnizeh/fingo/business/types/money"
	"github.com/garnizeh/fingo/business/types/name"
	"github.com/garnizeh/fingo/business/types/quantity"
	"github.com/google/uuid"
)

// TestGenerateNewProducts is a helper method for testing.
func TestGenerateNewProducts(n int, userID uuid.UUID) []NewProduct {
	newPrds := make([]NewProduct, n)

	idx := randomInt(10000)
	for i := range n {
		idx++

		np := NewProduct{
			Name:     name.MustParse(fmt.Sprintf("Name%d", idx)),
			Cost:     money.MustParse(float64(randomInt(500) + 1)),
			Quantity: quantity.MustParse(randomInt(50) + 1),
			UserID:   userID,
		}

		newPrds[i] = np
	}

	return newPrds
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

// TestGenerateSeedProducts is a helper method for testing.
func TestGenerateSeedProducts(ctx context.Context, n int, api ExtBusiness, userID uuid.UUID) ([]Product, error) {
	newPrds := TestGenerateNewProducts(n, userID)

	prds := make([]Product, len(newPrds))
	for i := range newPrds {
		prd, err := api.Create(ctx, &newPrds[i])
		if err != nil {
			return nil, fmt.Errorf("seeding product: idx: %d : %w", i, err)
		}

		prds[i] = prd
	}

	return prds, nil
}
