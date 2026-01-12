package main

import (
	"net/http"
	"secret-store/internal/api"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /secret/generate", api.GenerateSecretApi)

	http.ListenAndServe("[::]:8080", api.CORSMiddleware(mux))
}
