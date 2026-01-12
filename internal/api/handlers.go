package api

import (
	"encoding/json"
	"net/http"
	"secret-store/internal/secret"
	"strconv"
)

// GenerateSecretApi returns a Secret String for the given Parameters
func GenerateSecretApi(w http.ResponseWriter, r *http.Request) {
	lengthStr := r.URL.Query().Get("length")
	includeNumbersStr := r.URL.Query().Get("includeNumbers")
	includeSymbolsStr := r.URL.Query().Get("includeSymbols")

	length, err := strconv.Atoi(lengthStr)
	if err != nil || length < 8 {
		http.Error(w, "Invalid length", http.StatusBadRequest)
		return
	}

	includeNumbers := includeNumbersStr == "true"
	includeSymbols := includeSymbolsStr == "true"

	secret := secret.GenerateSecret(length, includeNumbers, includeSymbols)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"secret": secret,
	})
}
