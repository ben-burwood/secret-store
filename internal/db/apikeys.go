package db

import (
	"database/sql"
	"errors"
	"time"
)

type APIKey struct {
	ID        int64
	Name      string
	Key       string
	CreatedAt time.Time
}

func (s *Store) ListAPIKeys() ([]APIKey, error) {
	rows, err := s.Query("SELECT id, name, key, created_at FROM api_keys ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []APIKey
	for rows.Next() {
		var k APIKey
		if err := rows.Scan(&k.ID, &k.Name, &k.Key, &k.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, k)
	}
	return out, rows.Err()
}

func (s *Store) GetAPIKeyByID(id int64) (*APIKey, error) {
	row := s.QueryRow("SELECT id, name, key, created_at FROM api_keys WHERE id = ?", id)
	var k APIKey
	if err := row.Scan(&k.ID, &k.Name, &k.Key, &k.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &k, nil
}

func (s *Store) GetAPIKeyByToken(token string) (*APIKey, error) {
	row := s.QueryRow("SELECT id, name, key, created_at FROM api_keys WHERE key = ?", token)
	var k APIKey
	if err := row.Scan(&k.ID, &k.Name, &k.Key, &k.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &k, nil
}

// AllAPIKeyScopes returns a map of api_key_id -> []secret_id.
func (s *Store) AllAPIKeyScopes() (map[int64][]int64, error) {
	rows, err := s.Query("SELECT api_key_id, secret_id FROM api_key_secrets ORDER BY api_key_id, secret_id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[int64][]int64{}
	for rows.Next() {
		var keyID, secretID int64
		if err := rows.Scan(&keyID, &secretID); err != nil {
			return nil, err
		}
		out[keyID] = append(out[keyID], secretID)
	}
	return out, rows.Err()
}

// APIKeyScope returns the sorted list of secret_ids scoped to the given api_key_id.
func (s *Store) APIKeyScope(apiKeyID int64) ([]int64, error) {
	rows, err := s.Query("SELECT secret_id FROM api_key_secrets WHERE api_key_id = ? ORDER BY secret_id", apiKeyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (s *Store) UpdateAPIKeyToken(id int64, newToken string) error {
	_, err := s.Exec("UPDATE api_keys SET key = ? WHERE id = ?", newToken, id)
	return err
}

func (s *Store) DeleteAPIKey(id int64) error {
	_, err := s.Exec("DELETE FROM api_keys WHERE id = ?", id)
	return err
}

// ReplaceAPIKeyScopes replaces the api_key_id's scope rows with the supplied list.
// Caller is responsible for verifying the secret_ids exist.
func (s *Store) ReplaceAPIKeyScopes(apiKeyID int64, secretIDs []int64) error {
	tx, err := s.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck
	if _, err := tx.Exec("DELETE FROM api_key_secrets WHERE api_key_id = ?", apiKeyID); err != nil {
		return err
	}
	for _, sid := range secretIDs {
		if _, err := tx.Exec("INSERT INTO api_key_secrets (api_key_id, secret_id) VALUES (?, ?)", apiKeyID, sid); err != nil {
			return err
		}
	}
	return tx.Commit()
}
