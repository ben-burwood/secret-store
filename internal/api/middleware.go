package api

import (
	"net/http"
)

const validToken = "my-secret-token"

func TokenAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" || token != validToken {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized: missing or invalid token"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
