package session

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

const TTLSeconds int64 = 86400 // 24 hours, matches Python

type Store struct {
	mu       sync.Mutex
	sessions map[string]int64 // token -> expiry unix seconds
}

func New() *Store {
	return &Store{sessions: make(map[string]int64)}
}

func (s *Store) Create() string {
	now := time.Now().Unix()
	s.mu.Lock()
	defer s.mu.Unlock()
	for t, exp := range s.sessions {
		if exp < now {
			delete(s.sessions, t)
		}
	}
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		// crypto/rand should never fail; if it does, return empty string and the caller
		// will set a useless cookie. The next Validate will fail on empty token anyway.
		return ""
	}
	tok := hex.EncodeToString(b[:])
	s.sessions[tok] = now + TTLSeconds
	return tok
}

func (s *Store) Validate(token string) bool {
	if token == "" {
		return false
	}
	now := time.Now().Unix()
	s.mu.Lock()
	defer s.mu.Unlock()
	exp, ok := s.sessions[token]
	if !ok {
		return false
	}
	if exp < now {
		delete(s.sessions, token)
		return false
	}
	return true
}

func (s *Store) Delete(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, token)
}
