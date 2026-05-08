package handlers

import (
	"net/http"

	"github.com/ben-burwood/secret-store/internal/crypto"
	"github.com/ben-burwood/secret-store/internal/db"
	"github.com/ben-burwood/secret-store/internal/httpx"
)

type API struct {
	DB     *db.Store
	Crypto *crypto.Cipher
}

// GetSecret implements GET /api/secret?token=...&key=...
//
// Error precedence (matches the existing Python behaviour exactly):
//
//  1. token missing  → 401
//  2. key   missing  → 404
//  3. secret unknown → 404
//  4. token unknown  → 401
//  5. token's scope is non-empty AND secret not in it → 401
//  6. otherwise: 200 plain text decrypted value
func (h *API) GetSecret(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	key := r.URL.Query().Get("key")

	if token == "" {
		httpx.Text(w, http.StatusUnauthorized, "Unauthorized: missing or invalid token")
		return
	}
	if key == "" {
		httpx.Text(w, http.StatusNotFound, "Secret not found")
		return
	}

	secret, err := h.DB.GetSecretByKey(key)
	if err != nil {
		httpx.Text(w, http.StatusInternalServerError, "internal error")
		return
	}
	if secret == nil {
		httpx.Text(w, http.StatusNotFound, "Secret not found")
		return
	}

	apiKey, err := h.DB.GetAPIKeyByToken(token)
	if err != nil {
		httpx.Text(w, http.StatusInternalServerError, "internal error")
		return
	}
	if apiKey == nil {
		httpx.Text(w, http.StatusUnauthorized, "Unauthorized: missing or invalid token")
		return
	}

	scope, err := h.DB.APIKeyScope(apiKey.ID)
	if err != nil {
		httpx.Text(w, http.StatusInternalServerError, "internal error")
		return
	}
	if len(scope) > 0 {
		permitted := false
		for _, sid := range scope {
			if sid == secret.ID {
				permitted = true
				break
			}
		}
		if !permitted {
			httpx.Text(w, http.StatusUnauthorized, "Unauthorized: missing or invalid token")
			return
		}
	}

	plain, err := h.Crypto.Decrypt(secret.Value)
	if err != nil {
		httpx.Text(w, http.StatusInternalServerError, "decrypt failed")
		return
	}
	httpx.Text(w, http.StatusOK, plain)
}
