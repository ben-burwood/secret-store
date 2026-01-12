package store

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

const AuthKeyFile = storeDirectory + "/secrets.json.authkey"

func writeAuthKey(authKey string) error {
	return os.WriteFile(AuthKeyFile, []byte(authKey), 0644)
}

func GenerateAuthKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", fmt.Errorf("failed to generate key: %w", err)
	}
	writeAuthKey(base64.StdEncoding.EncodeToString(key))
	return base64.StdEncoding.EncodeToString(key), nil
}

func ReadAuthKey() (string, error) {
	if _, err := os.Stat(AuthKeyFile); os.IsNotExist(err) {
		// File Does Not Exist
		return "", nil
	}

	data, err := os.ReadFile(AuthKeyFile)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}
