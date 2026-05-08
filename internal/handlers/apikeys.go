package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ben-burwood/secret-store/internal/db"
	"github.com/ben-burwood/secret-store/internal/httpx"
)

type APIKeys struct {
	DB *db.Store
}

type apiKeyView struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Key       string  `json:"key"`
	CreatedAt string  `json:"created_at"`
	SecretIDs []int64 `json:"secret_ids"`
}

func newAPIKeyToken() string {
	var b [32]byte
	_, _ = rand.Read(b[:])
	return base64.StdEncoding.EncodeToString(b[:])
}

// sortedInt64s returns a non-nil sorted copy. Important so JSON encoding
// produces [] not null for empty scope sets.
func sortedInt64s(in []int64) []int64 {
	out := make([]int64, len(in))
	copy(out, in)
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

// dedupInt64s preserves nothing about order; the caller must sort if needed.
func dedupInt64s(in []int64) []int64 {
	seen := map[int64]struct{}{}
	out := make([]int64, 0, len(in))
	for _, v := range in {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

// parseSecretIDs returns (ids, true) when the JSON value is missing, null, or
// a list of integers. It returns (nil, false) when the value is present but is
// not a list of integers (matching Python's "secret_ids must be a list of integers").
func parseSecretIDs(raw json.RawMessage) ([]int64, bool) {
	if len(raw) == 0 {
		return nil, true
	}
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "null" {
		return nil, true
	}
	if !strings.HasPrefix(trimmed, "[") {
		return nil, false
	}
	var anyList []any
	if err := json.Unmarshal(raw, &anyList); err != nil {
		return nil, false
	}
	out := make([]int64, 0, len(anyList))
	for _, v := range anyList {
		f, ok := v.(float64)
		if !ok {
			return nil, false
		}
		if f != float64(int64(f)) {
			return nil, false
		}
		out = append(out, int64(f))
	}
	return out, true
}

func (h *APIKeys) List(w http.ResponseWriter, r *http.Request) {
	keys, err := h.DB.ListAPIKeys()
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	scopes, err := h.DB.AllAPIKeyScopes()
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	views := make([]apiKeyView, 0, len(keys))
	for _, k := range keys {
		ids := scopes[k.ID]
		if ids == nil {
			ids = []int64{}
		}
		views = append(views, apiKeyView{
			ID:        k.ID,
			Name:      k.Name,
			Key:       k.Key,
			CreatedAt: db.FormatISO(k.CreatedAt),
			SecretIDs: sortedInt64s(ids),
		})
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"api_keys": views})
}

func (h *APIKeys) Create(w http.ResponseWriter, r *http.Request) {
	var raw map[string]json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	var name string
	if rn, ok := raw["name"]; ok {
		_ = json.Unmarshal(rn, &name)
	}
	name = strings.TrimSpace(name)
	if name == "" {
		httpx.Error(w, http.StatusBadRequest, "name is required")
		return
	}
	rawSecretIDs := raw["secret_ids"]
	secretIDs, ok := parseSecretIDs(rawSecretIDs)
	if !ok {
		httpx.Error(w, http.StatusBadRequest, "secret_ids must be a list of integers")
		return
	}
	uniqueIDs := dedupInt64s(secretIDs)

	tx, err := h.DB.Begin()
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback() //nolint:errcheck

	if len(uniqueIDs) > 0 {
		exists, err := db.AllSecretsExist(tx, uniqueIDs)
		if err != nil {
			httpx.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		if !exists {
			httpx.Error(w, http.StatusBadRequest, "One or more secret_ids do not exist")
			return
		}
	}

	now := time.Now()
	keyToken := newAPIKeyToken()
	res, err := tx.Exec(
		"INSERT INTO api_keys (name, key, created_at) VALUES (?, ?, ?)",
		name, keyToken, db.FormatTime(now),
	)
	if err != nil {
		if db.IsUniqueViolation(err) {
			httpx.Error(w, http.StatusConflict, "Name '"+name+"' already exists")
			return
		}
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	id, _ := res.LastInsertId()
	for _, sid := range uniqueIDs {
		if _, err := tx.Exec(
			"INSERT INTO api_key_secrets (api_key_id, secret_id) VALUES (?, ?)",
			id, sid,
		); err != nil {
			httpx.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	if err := tx.Commit(); err != nil {
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpx.JSON(w, http.StatusCreated, apiKeyView{
		ID:        id,
		Name:      name,
		Key:       keyToken,
		CreatedAt: db.FormatISO(now),
		SecretIDs: sortedInt64s(uniqueIDs),
	})
}

func (h *APIKeys) Regenerate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "API key not found")
		return
	}
	existing, err := h.DB.GetAPIKeyByID(id)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if existing == nil {
		httpx.Error(w, http.StatusNotFound, "API key not found")
		return
	}
	newToken := newAPIKeyToken()
	if err := h.DB.UpdateAPIKeyToken(id, newToken); err != nil {
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	scope, err := h.DB.APIKeyScope(id)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if scope == nil {
		scope = []int64{}
	}
	httpx.JSON(w, http.StatusOK, apiKeyView{
		ID:        id,
		Name:      existing.Name,
		Key:       newToken,
		CreatedAt: db.FormatISO(existing.CreatedAt),
		SecretIDs: scope,
	})
}

func (h *APIKeys) UpdateScopes(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "API key not found")
		return
	}
	var raw map[string]json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	rawIDs, hasField := raw["secret_ids"]
	if !hasField {
		httpx.Error(w, http.StatusBadRequest, "secret_ids must be a list of integers")
		return
	}
	secretIDs, ok := parseSecretIDs(rawIDs)
	if !ok {
		httpx.Error(w, http.StatusBadRequest, "secret_ids must be a list of integers")
		return
	}
	uniqueIDs := dedupInt64s(secretIDs)

	existing, err := h.DB.GetAPIKeyByID(id)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if existing == nil {
		httpx.Error(w, http.StatusNotFound, "API key not found")
		return
	}

	tx, err := h.DB.Begin()
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback() //nolint:errcheck

	if len(uniqueIDs) > 0 {
		exists, err := db.AllSecretsExist(tx, uniqueIDs)
		if err != nil {
			httpx.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		if !exists {
			httpx.Error(w, http.StatusBadRequest, "One or more secret_ids do not exist")
			return
		}
	}
	if _, err := tx.Exec("DELETE FROM api_key_secrets WHERE api_key_id = ?", id); err != nil {
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	for _, sid := range uniqueIDs {
		if _, err := tx.Exec(
			"INSERT INTO api_key_secrets (api_key_id, secret_id) VALUES (?, ?)",
			id, sid,
		); err != nil {
			httpx.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	if err := tx.Commit(); err != nil {
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpx.JSON(w, http.StatusOK, apiKeyView{
		ID:        id,
		Name:      existing.Name,
		Key:       existing.Key,
		CreatedAt: db.FormatISO(existing.CreatedAt),
		SecretIDs: sortedInt64s(uniqueIDs),
	})
}

func (h *APIKeys) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		httpx.Empty(w, http.StatusNoContent)
		return
	}
	if err := h.DB.DeleteAPIKey(id); err != nil {
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.Empty(w, http.StatusNoContent)
}
