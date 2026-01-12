package secret

import (
	"math/rand"
)

const (
	Letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numbers = "0123456789"
	Symbols = "!@#$%^&*()-_=+[]{}|;:,.<>?/"
)

// GenerateSecret returns a Secret String for the given Parameters
func GenerateSecret(length int, includeNumbers, includeSymbols bool) string {
	charset := Letters
	if includeNumbers {
		charset += Numbers
	}
	if includeSymbols {
		charset += Symbols
	}

	secret := make([]byte, length)
	for i := range secret {
		secret[i] = charset[rand.Intn(len(charset))]
	}
	return string(secret)
}
