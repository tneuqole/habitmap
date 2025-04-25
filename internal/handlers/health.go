package handlers

import (
	"net/http"
)

func GetHealth(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("healthy")) //nolint:errcheck,gosec
	return nil
}
