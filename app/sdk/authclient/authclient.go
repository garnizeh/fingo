// Package authclient holds the authentication client relevant models and interfaces
package authclient

import (
	"context"
	"encoding/json"

	"github.com/garnizeh/fingo/app/sdk/auth"
	"github.com/google/uuid"
)

type Authenticator interface {
	Authenticate(ctx context.Context, authorization string) (AuthenticateResp, error)
	Authorize(ctx context.Context, auth *Authorize) error
	Close() error
}

// Authorize defines the information required to perform an authorization.
type Authorize struct {
	Rule   string
	Claims auth.Claims
	UserID uuid.UUID
}

// Decode implements the decoder interface.
func (a *Authorize) Decode(data []byte) error {
	return json.Unmarshal(data, a)
}

// AuthenticateResp defines the information that will be received on authenticate.
type AuthenticateResp struct {
	Claims auth.Claims
	UserID uuid.UUID
}

// Encode implements the encoder interface.
func (ar *AuthenticateResp) Encode() (data []byte, comtentType string, err error) {
	data, err = json.Marshal(ar)
	comtentType = "application/json"
	return
}
