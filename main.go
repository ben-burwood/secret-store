package main

import (
	"net/http"
	"secret-store/internal/api"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /secret/generate", api.GenerateSecretApi)
	mux.HandleFunc("GET /secrets", api.ListSecretsApi)
	mux.HandleFunc("POST /secrets/new", api.CreateSecretApi)
	mux.HandleFunc("DELETE /secrets/{id}", api.DeleteSecretApi)
	mux.HandleFunc("UPDATE /secrets/{id}", api.UpdateSecretApi)
	mux.HandleFunc("GET /auth/key", api.GetAuthKeyApi)
	mux.HandleFunc("GET /auth/key/generate", api.GenerateAuthKeyApi)

	http.ListenAndServe("[::]:8080", api.CORSMiddleware(mux))
}
