package creditcardbus

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/garnizeh/fingo/business/types/money"
	"github.com/garnizeh/fingo/business/types/name"
	"github.com/google/uuid"
)

// TestGenerateNewCreditCards is a helper method for testing.
func TestGenerateNewCreditCards(n int, userID uuid.UUID) []NewCreditCard {
	newCCs := make([]NewCreditCard, n)

	idx := rand.Intn(10000)
	for i := range newCCs {
		idx++

		ncc := NewCreditCard{
			UserID:         userID,
			Name:           name.MustParse(fmt.Sprintf("Card %d", idx)),
			Limit:          money.MustParse(float64(rand.Intn(10000) + 1000)),
			ClosingDay:     rand.Intn(28) + 1,
			DueDay:         rand.Intn(20) + 1,
			LastFourDigits: fmt.Sprintf("%04d", rand.Intn(10000)),
		}

		newCCs[i] = ncc
	}

	return newCCs
}

// TestGenerateSeedCreditCards is a helper method for testing.
func TestGenerateSeedCreditCards(ctx context.Context, n int, api ExtBusiness, userID uuid.UUID) ([]CreditCard, error) {
	newCCs := TestGenerateNewCreditCards(n, userID)

	ccs := make([]CreditCard, len(newCCs))
	for i, ncc := range newCCs {
		cc, err := api.Create(ctx, userID, ncc)
		if err != nil {
			return nil, fmt.Errorf("seeding credit card: idx: %d : %w", i, err)
		}

		ccs[i] = cc
	}

	return ccs, nil
}
