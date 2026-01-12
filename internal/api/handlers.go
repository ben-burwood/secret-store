package api

import (
	"net/http"
	"secret-store/internal/store"
)

func GetSecretApi(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	secret, err := store.GetSecret(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(secret.Value))
}
