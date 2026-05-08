package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ben-burwood/secret-store/internal/httpx"
	"github.com/ben-burwood/secret-store/internal/session"
)

type Auth struct {
	Sessions *session.Store
	User     string
	Pass     string
}

type loginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var body loginBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	if body.Username != a.User || body.Password != a.Pass {
		httpx.Error(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}
	tok := a.Sessions.Create()
	http.SetCookie(w, &http.Cookie{
		Name:     httpx.SessionCookie,
		Value:    tok,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	httpx.JSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (a *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(httpx.SessionCookie); err == nil {
		a.Sessions.Delete(c.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     httpx.SessionCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // emits Max-Age=0
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	httpx.JSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (a *Auth) Status(w http.ResponseWriter, r *http.Request) {
	httpx.Empty(w, http.StatusNoContent)
}
