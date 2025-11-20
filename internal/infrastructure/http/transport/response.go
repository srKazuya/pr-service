// Package transport provides HTTP transport utilities for writing JSON responses.
package transport

import (
	"encoding/json"
	"fmt"
	"net/http"
)


func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return fmt.Errorf("%v: %w", ErrEncode, err)
	}
	return nil
}
