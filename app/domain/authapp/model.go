package authapp

import "encoding/json"

type token struct {
	Token string `json:"token"`
}

// Encode implements the encoder interface.
func (t token) Encode() (data []byte, contentType string, err error) {
	data, err = json.Marshal(t)
	contentType = "application/json"
	return
}
