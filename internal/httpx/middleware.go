package httpx

import (
	"net/http"

	"github.com/ben-burwood/secret-store/internal/session"
)

func RequireSession(s *session.Store, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(SessionCookie)
		if err != nil || !s.Validate(c.Value) {
			Error(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}
