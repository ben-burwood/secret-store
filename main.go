package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/ben-burwood/secret-store/internal/config"
	"github.com/ben-burwood/secret-store/internal/crypto"
	"github.com/ben-burwood/secret-store/internal/db"
	"github.com/ben-burwood/secret-store/internal/handlers"
	"github.com/ben-burwood/secret-store/internal/httpx"
	"github.com/ben-burwood/secret-store/internal/session"
)

func main() {
	cfg := config.Load()

	cipher, err := crypto.New(cfg.EncryptionKey)
	if err != nil {
		log.Fatalf("ENCRYPTION_KEY invalid: %v", err)
	}

	store, err := db.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer store.Close()
	if err := store.Bootstrap(); err != nil {
		log.Fatalf("bootstrap db: %v", err)
	}

	sessions := session.New()

	authH := &handlers.Auth{Sessions: sessions, User: cfg.DashboardUser, Pass: cfg.DashboardPass}
	secretsH := &handlers.Secrets{DB: store, Crypto: cipher}
	apikeysH := &handlers.APIKeys{DB: store}
	apiH := &handlers.API{DB: store, Crypto: cipher}
	generateH := &handlers.Generate{}

	mux := http.NewServeMux()

	// Health & SPA-supporting public routes
	mux.HandleFunc("GET /health", handlers.Health)
	mux.HandleFunc("GET /web/secret/generate", generateH.Generate)
	mux.HandleFunc("GET /api/secret", apiH.GetSecret)

	// Auth endpoints
	mux.HandleFunc("POST /web/auth/login", authH.Login)
	mux.HandleFunc("POST /web/auth/logout", authH.Logout)
	mux.Handle("GET /web/auth/status", httpx.RequireSession(sessions, http.HandlerFunc(authH.Status)))

	// Secrets dashboard
	mux.Handle("GET /web/secrets", httpx.RequireSession(sessions, http.HandlerFunc(secretsH.List)))
	mux.Handle("POST /web/secrets/new", httpx.RequireSession(sessions, http.HandlerFunc(secretsH.Create)))
	mux.Handle("PATCH /web/secrets/{id}", httpx.RequireSession(sessions, http.HandlerFunc(secretsH.Update)))
	mux.Handle("DELETE /web/secrets/{id}", httpx.RequireSession(sessions, http.HandlerFunc(secretsH.Delete)))
	mux.Handle("POST /web/import", httpx.RequireSession(sessions, http.HandlerFunc(secretsH.Import)))
	mux.Handle("POST /web/restore", httpx.RequireSession(sessions, http.HandlerFunc(secretsH.Restore)))

	// API key management
	mux.Handle("GET /web/api/keys", httpx.RequireSession(sessions, http.HandlerFunc(apikeysH.List)))
	mux.Handle("POST /web/api/keys/new", httpx.RequireSession(sessions, http.HandlerFunc(apikeysH.Create)))
	mux.Handle("POST /web/api/keys/{id}/regenerate", httpx.RequireSession(sessions, http.HandlerFunc(apikeysH.Regenerate)))
	mux.Handle("PUT /web/api/keys/{id}/scopes", httpx.RequireSession(sessions, http.HandlerFunc(apikeysH.UpdateScopes)))
	mux.Handle("DELETE /web/api/keys/{id}", httpx.RequireSession(sessions, http.HandlerFunc(apikeysH.Delete)))

	// SPA fallback (registered last so it has lowest specificity)
	if cfg.StaticDir != "" {
		mux.Handle("/", httpx.SPA(cfg.StaticDir))
	}

	addr := "0.0.0.0:" + cfg.Port
	log.Printf("listening on %s (db=%s static=%s)", addr, cfg.DBPath, cfg.StaticDir)
	if err := http.ListenAndServe(addr, mux); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server: %v", err)
	}
}
