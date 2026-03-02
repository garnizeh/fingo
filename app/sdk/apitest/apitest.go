// Package apitest provides support for excuting api test logic.
package apitest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/mail"
	"testing"
	"time"

	"github.com/garnizeh/fingo/app/sdk/auth"
	"github.com/garnizeh/fingo/business/domain/userbus"
	"github.com/garnizeh/fingo/business/sdk/dbtest"
	"github.com/garnizeh/fingo/business/types/role"
	"github.com/golang-jwt/jwt/v4"
)

type testOption struct {
	skipMsg string
	skip    bool
}

type OptionFunc func(*testOption)

// WithSkip can be used to skip running a test.
func WithSkip(skip bool, msg string) OptionFunc {
	return func(to *testOption) {
		to.skip = skip
		to.skipMsg = msg
	}
}

// Test contains functions for executing an api test.
type Test struct {
	DB   *dbtest.Database
	Auth *auth.Auth
	mux  http.Handler
}

// Run performs the actual test logic based on the table data.
func (at *Test) Run(t *testing.T, table []Table, testName string, options ...OptionFunc) {
	to := new(testOption)
	for _, f := range options {
		f(to)
	}

	if to.skip {
		t.Skipf("%v: %v", testName, to.skipMsg)
	}

	for i := range table {
		tt := &table[i]
		f := func(t *testing.T) {
			var body io.Reader = http.NoBody
			w := httptest.NewRecorder()

			if tt.Input != nil {
				d, err := json.Marshal(tt.Input)
				if err != nil {
					t.Fatalf("Should be able to marshal the model : %s", err)
				}

				body = bytes.NewBuffer(d)
			}

			r := httptest.NewRequest(tt.Method, tt.URL, body)

			r.Header.Set("Authorization", "Bearer "+tt.Token)
			at.mux.ServeHTTP(w, r)

			if w.Code != tt.StatusCode {
				t.Fatalf("%s: Should receive a status code of %d for the response : %d", tt.Name, tt.StatusCode, w.Code)
			}

			if tt.StatusCode == http.StatusNoContent {
				return
			}

			if err := json.Unmarshal(w.Body.Bytes(), tt.GotResp); err != nil {
				t.Fatalf("Should be able to unmarshal the response : %s", err)
			}

			diff := tt.CmpFunc(tt.GotResp, tt.ExpResp)
			if diff != "" {
				t.Log("DIFF")
				t.Logf("%s", diff)
				t.Log("GOT")
				t.Logf("%#v", tt.GotResp)
				t.Log("EXP")
				t.Logf("%#v", tt.ExpResp)
				t.Fatalf("Should get the expected response")
			}
		}

		t.Run(testName+"-"+tt.Name, f)
	}
}

// =============================================================================

// Token generates an authenticated token for a user.
func Token(userBus userbus.ExtBusiness, ath *auth.Auth, email string) string {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return ""
	}

	dbUsr, err := userBus.QueryByEmail(context.Background(), *addr)
	if err != nil {
		return ""
	}

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   dbUsr.ID.String(),
			Issuer:    ath.Issuer(),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: role.ParseToString(dbUsr.Roles),
	}

	token, err := ath.GenerateToken(kid, &claims)
	if err != nil {
		return ""
	}

	return token
}
