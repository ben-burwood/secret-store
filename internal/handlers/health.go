package handlers

import (
	"net/http"

	"github.com/ben-burwood/secret-store/internal/httpx"
)

func Health(w http.ResponseWriter, r *http.Request) {
	httpx.Text(w, http.StatusOK, "ok")
}
