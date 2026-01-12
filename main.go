package main

import (
	"net/http"
	"secret-store/internal/api"
	"secret-store/internal/store"
	"secret-store/internal/web"
)

func main() {
	store.InitialiseStore()

	webMux := http.NewServeMux()
	webMux.HandleFunc("GET /web/secret/generate", web.GenerateSecretWeb)
	webMux.HandleFunc("GET /web/secrets", web.ListSecretsWeb)
	webMux.HandleFunc("POST /web/secrets/new", web.CreateSecretWeb)
	webMux.HandleFunc("DELETE /web/secrets/{id}", web.DeleteSecretWeb)
	webMux.HandleFunc("PATCH /web/secrets/{id}", web.UpdateSecretWeb)
	webMux.HandleFunc("GET /web/auth/key", web.GetAuthKeyWeb)
	webMux.HandleFunc("GET /web/auth/key/generate", web.GenerateAuthKeyWeb)
	// Serve Static Frontend
	webMux.Handle("/", http.FileServer(http.Dir("./frontend/dist")))

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("GET /api/secret", api.GetSecretApi)

	// Start web server on 8080 in a goroutine
	go http.ListenAndServe("[::]:8080", web.CORSMiddleware(webMux))

	// Start API server on 8000 (main goroutine)
	http.ListenAndServe("[::]:8000", api.TokenAuthMiddleware(apiMux))
}
