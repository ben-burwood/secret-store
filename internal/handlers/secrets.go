package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ben-burwood/secret-store/internal/crypto"
	"github.com/ben-burwood/secret-store/internal/db"
	"github.com/ben-burwood/secret-store/internal/httpx"
)

type Secrets struct {
	DB     *db.Store
	Crypto *crypto.Cipher
}

type secretView struct {
	ID        int64   `json:"id"`
	Key       string  `json:"key"`
	Value     string  `json:"value"`
	Tag       *string `json:"tag"`
	CreatedAt string  `json:"created_at"`
}

func toSecretView(s db.Secret, plain string) secretView {
	v := secretView{
		ID:        s.ID,
		Key:       s.Key,
		Value:     plain,
		CreatedAt: db.FormatISO(s.CreatedAt),
	}
	if s.Tag.Valid {
		t := s.Tag.String
		v.Tag = &t
	}
	return v
}

func normalizeTag(value any) *string {
	s, ok := value.(string)
	if !ok {
		return nil
	}
	n := strings.ReplaceAll(strings.ToLower(strings.TrimSpace(s)), " ", "_")
	if n == "" {
		return nil
	}
	return &n
}

func (h *Secrets) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.ListSecrets()
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	views := make([]secretView, 0, len(rows))
	for _, s := range rows {
		plain, err := h.Crypto.Decrypt(s.Value)
		if err != nil {
			httpx.Error(w, http.StatusInternalServerError, "decrypt failed")
			return
		}
		views = append(views, toSecretView(s, plain))
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"secrets": views})
}

func (h *Secrets) Create(w http.ResponseWriter, r *http.Request) {
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	keyRaw, _ := body["key"].(string)
	key := strings.TrimSpace(keyRaw)
	if key == "" {
		httpx.Error(w, http.StatusBadRequest, "key is required")
		return
	}
	valRaw, _ := body["value"].(string)
	enc, err := h.Crypto.Encrypt(valRaw)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "encrypt failed")
		return
	}
	tag := normalizeTag(body["tag"])
	if _, err := h.DB.CreateSecret(key, enc, tag, time.Now()); err != nil {
		if db.IsUniqueViolation(err) {
			httpx.Error(w, http.StatusConflict, "Key '"+key+"' already exists")
			return
		}
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.Empty(w, http.StatusCreated)
}

func (h *Secrets) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid id")
		return
	}
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	existing, err := h.DB.GetSecretByID(id)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if existing == nil {
		httpx.Error(w, http.StatusNotFound, "Secret not found")
		return
	}

	updates := map[string]any{}
	var newKey string
	if rawKey, has := body["key"]; has {
		ks, _ := rawKey.(string)
		k := strings.TrimSpace(ks)
		if k == "" {
			httpx.Error(w, http.StatusBadRequest, "key is required")
			return
		}
		if k != existing.Key {
			updates["key"] = k
			newKey = k
		}
	}
	if rawVal, has := body["value"]; has {
		vs, _ := rawVal.(string)
		enc, err := h.Crypto.Encrypt(vs)
		if err != nil {
			httpx.Error(w, http.StatusInternalServerError, "encrypt failed")
			return
		}
		updates["value"] = enc
	}
	if _, has := body["tag"]; has {
		t := normalizeTag(body["tag"])
		if t == nil {
			updates["tag"] = nil
		} else {
			updates["tag"] = *t
		}
	}
	updates["created_at"] = db.FormatTime(time.Now())

	if err := h.DB.UpdateSecretFields(id, updates); err != nil {
		if db.IsUniqueViolation(err) {
			httpx.Error(w, http.StatusConflict, "Key '"+newKey+"' already exists")
			return
		}
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.Empty(w, http.StatusNoContent)
}

func (h *Secrets) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		httpx.Empty(w, http.StatusNoContent)
		return
	}
	if err := h.DB.DeleteSecret(id); err != nil {
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.Empty(w, http.StatusNoContent)
}

type secretImportItem struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Tag       any    `json:"tag"`
	CreatedAt string `json:"created_at"`
}

func readSecretsForm(r *http.Request) ([]secretImportItem, string) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return nil, "Missing 'secrets' field"
	}
	raw := r.FormValue("secrets")
	if raw == "" {
		return nil, "Missing 'secrets' field"
	}
	var asList []secretImportItem
	if err := json.Unmarshal([]byte(raw), &asList); err == nil {
		return asList, ""
	}
	var asObj struct {
		Secrets *[]secretImportItem `json:"secrets"`
	}
	if err := json.Unmarshal([]byte(raw), &asObj); err != nil {
		return nil, "Invalid JSON in 'secrets'"
	}
	if asObj.Secrets == nil {
		return nil, "'secrets' must be a list"
	}
	return *asObj.Secrets, ""
}

type dedupedItem struct {
	Key       string
	Value     string
	Tag       *string
	CreatedAt time.Time
}

func dedupItems(items []secretImportItem, defaultCreatedAt time.Time, preserveCreatedAt bool) []dedupedItem {
	byKey := map[string]int{}
	out := []dedupedItem{}
	for _, item := range items {
		k := strings.TrimSpace(item.Key)
		if k == "" {
			continue
		}
		ts := defaultCreatedAt
		if preserveCreatedAt {
			if parsed, ok := db.ParseFlexibleISO(item.CreatedAt); ok {
				ts = parsed
			}
		}
		entry := dedupedItem{
			Key:       k,
			Value:     item.Value,
			Tag:       normalizeTag(item.Tag),
			CreatedAt: ts,
		}
		if idx, exists := byKey[k]; exists {
			out[idx] = entry
			continue
		}
		byKey[k] = len(out)
		out = append(out, entry)
	}
	return out
}

func (h *Secrets) Import(w http.ResponseWriter, r *http.Request) {
	items, errMsg := readSecretsForm(r)
	if errMsg != "" {
		httpx.Error(w, http.StatusBadRequest, errMsg)
		return
	}
	now := time.Now()
	deduped := dedupItems(items, now, false)

	tx, err := h.DB.Begin()
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback() //nolint:errcheck

	for _, item := range deduped {
		enc, err := h.Crypto.Encrypt(item.Value)
		if err != nil {
			httpx.Error(w, http.StatusInternalServerError, "encrypt failed")
			return
		}
		var tagVal sql.NullString
		if item.Tag != nil {
			tagVal = sql.NullString{String: *item.Tag, Valid: true}
		}
		if _, err := tx.Exec(`
			INSERT INTO secrets (key, value, tag, created_at)
			VALUES (?, ?, ?, ?)
			ON CONFLICT(key) DO UPDATE SET
				value = excluded.value,
				tag = excluded.tag,
				created_at = excluded.created_at
		`, item.Key, enc, tagVal, db.FormatTime(item.CreatedAt)); err != nil {
			httpx.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	if err := tx.Commit(); err != nil {
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.Empty(w, http.StatusNoContent)
}

func (h *Secrets) Restore(w http.ResponseWriter, r *http.Request) {
	items, errMsg := readSecretsForm(r)
	if errMsg != "" {
		httpx.Error(w, http.StatusBadRequest, errMsg)
		return
	}
	deduped := dedupItems(items, time.Now(), true)

	tx, err := h.DB.Begin()
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback() //nolint:errcheck

	if _, err := tx.Exec("DELETE FROM secrets"); err != nil {
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	for _, item := range deduped {
		enc, err := h.Crypto.Encrypt(item.Value)
		if err != nil {
			httpx.Error(w, http.StatusInternalServerError, "encrypt failed")
			return
		}
		var tagVal sql.NullString
		if item.Tag != nil {
			tagVal = sql.NullString{String: *item.Tag, Valid: true}
		}
		if _, err := tx.Exec(
			"INSERT INTO secrets (key, value, tag, created_at) VALUES (?, ?, ?, ?)",
			item.Key, enc, tagVal, db.FormatTime(item.CreatedAt),
		); err != nil {
			httpx.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	if err := tx.Commit(); err != nil {
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.Empty(w, http.StatusNoContent)
}
