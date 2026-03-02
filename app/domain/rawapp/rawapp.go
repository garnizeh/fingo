// Package rawapp provides an example of using a raw handler.
package rawapp

import (
	"encoding/json"
	"net/http"
)

func rawHandler(w http.ResponseWriter, r *http.Request) {
	status := struct {
		Status string
	}{
		Status: "RAW ENDPOINT",
	}

	if err := json.NewEncoder(w).Encode(status); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
