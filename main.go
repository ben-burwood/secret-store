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

	http.ListenAndServe("[::]:8080", api.CORSMiddleware(mux))
}
