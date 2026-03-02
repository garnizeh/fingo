package auditbus

import (
	"context"
	cryptorand "crypto/rand"
	"fmt"
	"math/big"

	"github.com/garnizeh/fingo/business/types/domain"
	"github.com/garnizeh/fingo/business/types/name"
	"github.com/google/uuid"
)

// TestNewAudits is a helper method for testing.
func TestNewAudits(n int, actorID uuid.UUID, objDomain domain.Domain, action string) []NewAudit {
	newAudits := make([]NewAudit, n)

	idx := randomInt(10000)
	for i := range n {
		idx++

		na := NewAudit{
			ObjID:     uuid.New(),
			ObjDomain: objDomain,
			ObjName:   name.MustParse(fmt.Sprintf("ObjName%d", idx)),
			ActorID:   actorID,
			Action:    action,
			Data:      struct{ Name string }{Name: fmt.Sprintf("Name%d", idx)},
			Message:   fmt.Sprintf("Message%d", idx),
		}

		newAudits[i] = na
	}

	return newAudits
}

// TestSeedAudits is a helper method for testing.
func TestSeedAudits(ctx context.Context, n int, actorID uuid.UUID, objDomain domain.Domain, action string, api ExtBusiness) ([]Audit, error) {
	newAudits := TestNewAudits(n, actorID, objDomain, action)

	audits := make([]Audit, len(newAudits))
	for i, na := range newAudits {
		adt, err := api.Create(ctx, &na)
		if err != nil {
			return nil, fmt.Errorf("seeding audit: idx: %d : %w", i, err)
		}

		audits[i] = adt
	}

	return audits, nil
}

func randomInt(upper int) int {
	if upper <= 0 {
		return 0
	}
	n, err := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(upper)))
	if err != nil {
		panic(err)
	}
	return int(n.Int64())
}
