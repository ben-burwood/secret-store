package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

const (
	SQLiteTimeLayout = "2006-01-02 15:04:05.000000"
	ISOOutputLayout  = "2006-01-02T15:04:05.000000"
)

type Store struct {
	*sql.DB
}

func Open(path string) (*Store, error) {
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("ensure dir %s: %w", dir, err)
		}
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	for _, p := range []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA foreign_keys=ON",
		"PRAGMA busy_timeout=5000",
	} {
		if _, err := db.Exec(p); err != nil {
			_ = db.Close()
			return nil, fmt.Errorf("set pragma %q: %w", p, err)
		}
	}
	return &Store{DB: db}, nil
}

func (s *Store) Bootstrap() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS secrets (
			id INTEGER NOT NULL PRIMARY KEY,
			key VARCHAR NOT NULL,
			value VARCHAR NOT NULL,
			tag VARCHAR,
			created_at DATETIME NOT NULL,
			UNIQUE (key)
		)`,
		`CREATE TABLE IF NOT EXISTS api_keys (
			id INTEGER NOT NULL PRIMARY KEY,
			name VARCHAR NOT NULL,
			key VARCHAR NOT NULL,
			created_at DATETIME NOT NULL,
			UNIQUE (name),
			UNIQUE (key)
		)`,
		`CREATE TABLE IF NOT EXISTS api_key_secrets (
			api_key_id INTEGER NOT NULL,
			secret_id INTEGER NOT NULL,
			PRIMARY KEY (api_key_id, secret_id),
			FOREIGN KEY(api_key_id) REFERENCES api_keys (id) ON DELETE CASCADE,
			FOREIGN KEY(secret_id) REFERENCES secrets (id) ON DELETE CASCADE
		)`,
	}
	for _, q := range stmts {
		if _, err := s.Exec(q); err != nil {
			return fmt.Errorf("bootstrap: %w", err)
		}
	}
	return nil
}

// IsUniqueViolation reports whether err is a SQLite UNIQUE constraint failure.
func IsUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}

// FormatTime emits the SQLAlchemy-compatible datetime string used in the schema.
// Reads come back via direct time.Time Scan; modernc.org/sqlite parses the stored
// string transparently, so no read-side helper is needed.
func FormatTime(t time.Time) string {
	return t.Format(SQLiteTimeLayout)
}

// FormatISO emits a Python datetime.isoformat()-equivalent string for naive datetimes.
func FormatISO(t time.Time) string {
	return t.Format(ISOOutputLayout)
}

// ParseFlexibleISO accepts the various ISO formats clients may send for
// created_at on /web/restore, mirroring Python's
// datetime.fromisoformat(value.replace("Z", "+00:00")).
func ParseFlexibleISO(raw string) (time.Time, bool) {
	if raw == "" {
		return time.Time{}, false
	}
	normalized := raw
	if n := len(normalized); n > 0 && normalized[n-1] == 'Z' {
		normalized = normalized[:n-1] + "+00:00"
	}
	layouts := []string{
		ISOOutputLayout,
		"2006-01-02T15:04:05",
		time.RFC3339Nano,
		time.RFC3339,
	}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, normalized, time.Local); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}
