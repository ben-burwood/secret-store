package api

import (
	"encoding/json"
	"net/http"
	"secret-store/internal/auth"
	"secret-store/internal/secret"
	"secret-store/internal/store"
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

func GetAuthKeyApi(w http.ResponseWriter, r *http.Request) {
	// TODO - Change this to return previously generated key or none
	key, err := auth.GenerateAuthKey()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"key": key,
	})
}

func GenerateAuthKeyApi(w http.ResponseWriter, r *http.Request) {
	key, err := auth.GenerateAuthKey()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"key": key,
	})
}

func ListSecretsApi(w http.ResponseWriter, r *http.Request) {
	secrets := store.ListSecrets()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]store.Secret{
		"secrets": secrets,
	})
}

func CreateSecretApi(w http.ResponseWriter, r *http.Request) {
	var secret store.Secret

	err := json.NewDecoder(r.Body).Decode(&secret)
	if err != nil {
		http.Error(w, "Invalid secret", http.StatusBadRequest)
		return
	}

	_, err = store.CreateSecret(secret.Key, secret.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeleteSecretApi(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = store.DeleteSecret(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateSecretApi(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var secret store.Secret

	err = json.NewDecoder(r.Body).Decode(&secret)
	if err != nil {
		http.Error(w, "Invalid secret", http.StatusBadRequest)
		return
	}

	err = store.UpdateSecret(id, secret.Key, secret.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
