package db

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

type Secret struct {
	ID        int64
	Key       string
	Value     string // raw stored ciphertext (Fernet-encrypted)
	Tag       sql.NullString
	CreatedAt time.Time
}

func (s *Store) ListSecrets() ([]Secret, error) {
	rows, err := s.Query("SELECT id, key, value, tag, created_at FROM secrets ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Secret
	for rows.Next() {
		sec, err := scanSecret(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, sec)
	}
	return out, rows.Err()
}

func (s *Store) GetSecretByKey(key string) (*Secret, error) {
	row := s.QueryRow("SELECT id, key, value, tag, created_at FROM secrets WHERE key = ?", key)
	sec, err := scanSecretRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &sec, nil
}

func (s *Store) GetSecretByID(id int64) (*Secret, error) {
	row := s.QueryRow("SELECT id, key, value, tag, created_at FROM secrets WHERE id = ?", id)
	sec, err := scanSecretRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &sec, nil
}

func (s *Store) CreateSecret(key, encryptedValue string, tag *string, createdAt time.Time) (int64, error) {
	res, err := s.Exec(
		"INSERT INTO secrets (key, value, tag, created_at) VALUES (?, ?, ?, ?)",
		key, encryptedValue, nullableTag(tag), FormatTime(createdAt),
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateSecretFields applies a partial UPDATE over the named columns.
// Values for "created_at" must already be SQLite-formatted strings.
// A nil value for "tag" produces a NULL in the database.
func (s *Store) UpdateSecretFields(id int64, fields map[string]any) error {
	if len(fields) == 0 {
		return nil
	}
	cols := make([]string, 0, len(fields))
	for c := range fields {
		cols = append(cols, c)
	}
	sort.Strings(cols)
	args := make([]any, 0, len(fields)+1)
	parts := make([]string, 0, len(fields))
	for _, c := range cols {
		parts = append(parts, c+" = ?")
		args = append(args, fields[c])
	}
	args = append(args, id)
	q := "UPDATE secrets SET " + strings.Join(parts, ", ") + " WHERE id = ?"
	_, err := s.Exec(q, args...)
	return err
}

func (s *Store) DeleteSecret(id int64) error {
	_, err := s.Exec("DELETE FROM secrets WHERE id = ?", id)
	return err
}

// AllSecretsExist reports whether every id in ids exists in the secrets table.
// `ids` must contain unique values; callers should de-dup beforehand.
func AllSecretsExist(q Queryer, ids []int64) (bool, error) {
	if len(ids) == 0 {
		return true, nil
	}
	placeholders := strings.Repeat("?,", len(ids))
	placeholders = placeholders[:len(placeholders)-1]
	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	query := fmt.Sprintf("SELECT COUNT(*) FROM secrets WHERE id IN (%s)", placeholders)
	var n int
	if err := q.QueryRow(query, args...).Scan(&n); err != nil {
		return false, err
	}
	return n == len(ids), nil
}

// Queryer is satisfied by both *sql.DB and *sql.Tx.
type Queryer interface {
	QueryRow(query string, args ...any) *sql.Row
}

func scanSecret(rows *sql.Rows) (Secret, error) {
	var sec Secret
	if err := rows.Scan(&sec.ID, &sec.Key, &sec.Value, &sec.Tag, &sec.CreatedAt); err != nil {
		return Secret{}, err
	}
	return sec, nil
}

func scanSecretRow(row *sql.Row) (Secret, error) {
	var sec Secret
	if err := row.Scan(&sec.ID, &sec.Key, &sec.Value, &sec.Tag, &sec.CreatedAt); err != nil {
		return Secret{}, err
	}
	return sec, nil
}

func nullableTag(tag *string) sql.NullString {
	if tag == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *tag, Valid: true}
}
