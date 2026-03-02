package creditcardbus

import (
	"context"
	cryptorand "crypto/rand"
	"fmt"
	"math/big"

	"github.com/garnizeh/fingo/business/types/money"
	"github.com/garnizeh/fingo/business/types/name"
	"github.com/google/uuid"
)

// TestGenerateNewCreditCards is a helper method for testing.
func TestGenerateNewCreditCards(n int, userID uuid.UUID) []NewCreditCard {
	newCCs := make([]NewCreditCard, n)

	idx := randomInt(10000)
	for i := range newCCs {
		idx++

		ncc := NewCreditCard{
			CreditCardIdentity: CreditCardIdentity{
				Name:           name.MustParse(fmt.Sprintf("Card %d", idx)),
				LastFourDigits: fmt.Sprintf("%04d", randomInt(10000)),
			},
			UserID:     userID,
			Limit:      money.MustParse(float64(randomInt(10000) + 1000)),
			ClosingDay: randomInt(28) + 1,
			DueDay:     randomInt(20) + 1,
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
