package main

import (
	"net/http"
	"secret-store/internal/web"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /web/secret/generate", web.GenerateSecretApi)
	mux.HandleFunc("GET /web/secrets", web.ListSecretsApi)
	mux.HandleFunc("POST /web/secrets/new", web.CreateSecretApi)
	mux.HandleFunc("DELETE /web/secrets/{id}", web.DeleteSecretApi)
	mux.HandleFunc("UPDATE /web/secrets/{id}", web.UpdateSecretApi)
	mux.HandleFunc("GET /web/auth/key", web.GetAuthKeyApi)
	mux.HandleFunc("GET /web/auth/key/generate", web.GenerateAuthKeyApi)

	http.ListenAndServe("[::]:8080", web.CORSMiddleware(mux))
}
