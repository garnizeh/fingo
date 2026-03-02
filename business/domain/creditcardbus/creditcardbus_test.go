package creditcardbus_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/garnizeh/fingo/business/domain/creditcardbus"
	"github.com/garnizeh/fingo/business/domain/userbus"
	"github.com/garnizeh/fingo/business/sdk/dbtest"
	"github.com/garnizeh/fingo/business/sdk/page"
	"github.com/garnizeh/fingo/business/types/money"
	"github.com/garnizeh/fingo/business/types/name"
	"github.com/garnizeh/fingo/business/types/role"
	"github.com/google/uuid"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	ctx, api, userAID, _ := newTestBus(t)

	ncc := creditcardbus.NewCreditCard{
		CreditCardIdentity: creditcardbus.CreditCardIdentity{
			Name:           name.MustParse("Primary Card"),
			LastFourDigits: "1234",
		},
		UserID:     userAID,
		Limit:      money.MustParse(1200),
		ClosingDay: 10,
		DueDay:     20,
	}

	cc, err := api.Create(ctx, userAID, ncc)
	if err != nil {
		t.Fatalf("create card: %v", err)
	}

	if cc.ID == uuid.Nil {
		t.Fatal("expected non-zero credit card id")
	}
	if cc.UserID != userAID {
		t.Fatalf("expected user id %s, got %s", userAID, cc.UserID)
	}
	if cc.Name != ncc.Name {
		t.Fatalf("expected name %s, got %s", ncc.Name, cc.Name)
	}
	if cc.Limit != ncc.Limit {
		t.Fatalf("expected limit %s, got %s", ncc.Limit, cc.Limit)
	}
	if !cc.Enabled {
		t.Fatal("expected card enabled by default")
	}

	invalid := creditcardbus.NewCreditCard{
		CreditCardIdentity: creditcardbus.CreditCardIdentity{
			Name:           name.MustParse("Invalid Limit Card"),
			LastFourDigits: "9999",
		},
		UserID:     userAID,
		Limit:      money.Money{},
		ClosingDay: 10,
		DueDay:     20,
	}

	_, err = api.Create(ctx, userAID, invalid)
	if !errors.Is(err, creditcardbus.ErrCardLimit) {
		t.Fatalf("expected ErrCardLimit, got %v", err)
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	ctx, api, userAID, _ := newTestBus(t)

	created := mustCreateCard(t, ctx, api, userAID, "Old Card", 1000, "1000")

	ucc := creditcardbus.UpdateCreditCard{
		Name:  dbtest.NamePointer("Updated Card"),
		Limit: dbtest.MoneyPointer(1500),
	}

	updated, err := api.Update(ctx, userAID, &created, ucc)
	if err != nil {
		t.Fatalf("update card: %v", err)
	}

	if updated.Name.String() != "Updated Card" {
		t.Fatalf("expected updated name, got %s", updated.Name)
	}
	if updated.Limit.Value() != 1500 {
		t.Fatalf("expected updated limit 1500, got %.2f", updated.Limit.Value())
	}
	if !updated.DateUpdated.After(created.DateUpdated) {
		t.Fatalf("expected date_updated to be after %v, got %v", created.DateUpdated, updated.DateUpdated)
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	ctx, api, userAID, _ := newTestBus(t)

	created := mustCreateCard(t, ctx, api, userAID, "Delete Card", 900, "2222")

	if err := api.Delete(ctx, userAID, &created); err != nil {
		t.Fatalf("delete card: %v", err)
	}

	_, err := api.QueryByID(ctx, userAID, created.ID)
	if !errors.Is(err, creditcardbus.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestQuery(t *testing.T) {
	t.Parallel()

	ctx, api, userAID, _ := newTestBus(t)

	mustCreateCard(t, ctx, api, userAID, "Card A1", 1000, "1111")
	mustCreateCard(t, ctx, api, userAID, "Card A2", 1100, "2222")
	mustCreateCard(t, ctx, api, userAID, "Card A3", 1200, "3333")

	filter := creditcardbus.QueryFilter{UserID: &userAID}

	cards, err := api.Query(ctx, userAID, filter, creditcardbus.DefaultOrderBy, page.MustParse("1", "10"))
	if err != nil {
		t.Fatalf("query cards: %v", err)
	}

	if len(cards) != 3 {
		t.Fatalf("expected 3 cards, got %d", len(cards))
	}
}

func TestQuery_UserIsolation(t *testing.T) {
	t.Parallel()

	ctx, api, userAID, userBID := newTestBus(t)

	mustCreateCard(t, ctx, api, userAID, "Card A", 1200, "1111")
	bCard := mustCreateCard(t, ctx, api, userBID, "Card B", 1300, "2222")

	filter := creditcardbus.QueryFilter{UserID: &userAID}

	cards, err := api.Query(ctx, userAID, filter, creditcardbus.DefaultOrderBy, page.MustParse("1", "10"))
	if err != nil {
		t.Fatalf("query cards by user A: %v", err)
	}

	for _, card := range cards {
		if card.UserID != userAID {
			t.Fatalf("found card from another user: expected %s, got %s", userAID, card.UserID)
		}
		if card.ID == bCard.ID {
			t.Fatalf("user B card %s leaked into user A query", bCard.ID)
		}
	}
}

func TestQueryByID(t *testing.T) {
	t.Parallel()

	ctx, api, userAID, _ := newTestBus(t)

	created := mustCreateCard(t, ctx, api, userAID, "ByID Card", 1400, "4444")

	got, err := api.QueryByID(ctx, userAID, created.ID)
	if err != nil {
		t.Fatalf("query by id: %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("expected id %s, got %s", created.ID, got.ID)
	}

	_, err = api.QueryByID(ctx, userAID, uuid.New())
	if !errors.Is(err, creditcardbus.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for unknown id, got %v", err)
	}
}

func newTestBus(t *testing.T) (ctx context.Context, api creditcardbus.ExtBusiness, userAID, userBID uuid.UUID) {
	t.Helper()

	if os.Getenv("TEST_DB_URL") == "" {
		t.Skip("skipping db test: TEST_DB_URL is not set")
	}

	db := dbtest.New(t, t.Name())
	ctx = context.Background()

	users, err := userbus.TestSeedUsers(ctx, 2, role.User, db.BusDomain.User)
	if err != nil {
		t.Fatalf("seeding users: %v", err)
	}

	api = db.BusDomain.CreditCard

	return ctx, api, users[0].ID, users[1].ID
}
func mustCreateCard(t *testing.T, ctx context.Context, api creditcardbus.ExtBusiness, userID uuid.UUID, cardName string, limit float64, lastFour string) creditcardbus.CreditCard {
	t.Helper()

	ncc := creditcardbus.NewCreditCard{
		CreditCardIdentity: creditcardbus.CreditCardIdentity{
			Name:           name.MustParse(cardName),
			LastFourDigits: lastFour,
		},
		UserID:     userID,
		Limit:      money.MustParse(limit),
		ClosingDay: 10,
		DueDay:     20,
	}

	cc, err := api.Create(ctx, userID, ncc)
	if err != nil {
		t.Fatalf("create card %s: %v", cardName, err)
	}

	return cc
}
