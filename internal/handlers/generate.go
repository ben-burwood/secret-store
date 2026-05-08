package handlers

import (
	"crypto/rand"
	"math/big"
	"net/http"
	"strconv"

	"github.com/ben-burwood/secret-store/internal/httpx"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits  = "0123456789"
	symbols = "!@#$%^&*()-_=+[]{}|;:,.<>?/"
)

type Generate struct{}

func truthy(v string) bool {
	return v == "true" || v == "True" || v == "TRUE"
}

func (g *Generate) Generate(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	rawLen := q.Get("length")
	length, err := strconv.Atoi(rawLen)
	if err != nil || length < 8 {
		httpx.Error(w, http.StatusBadRequest, "Invalid length")
		return
	}

	charset := letters
	if truthy(q.Get("includeNumbers")) {
		charset += digits
	}
	if truthy(q.Get("includeSymbols")) {
		charset += symbols
	}

	max := big.NewInt(int64(len(charset)))
	out := make([]byte, length)
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			httpx.Error(w, http.StatusInternalServerError, "rand failed")
			return
		}
		out[i] = charset[n.Int64()]
	}

	httpx.JSON(w, http.StatusOK, map[string]string{"secret": string(out)})
}
