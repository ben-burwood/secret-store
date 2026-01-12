package store

import (
	"errors"
	"math/rand"
	"time"
)

var ErrSecretNotFound = errors.New("secret not found")

type Secret struct {
	ID        int       `json:"id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

func ListSecrets() ([]Secret, error) {
	return loadSecrets()
}

func GetSecret(key string) (Secret, error) {
	secrets, err := ListSecrets()
	if err != nil {
		return Secret{}, err
	}
	for _, secret := range secrets {
		if secret.Key == key {
			return secret, nil
		}
	}
	return Secret{}, ErrSecretNotFound
}

func CreateSecret(key string, value string) (Secret, error) {
	secret := Secret{
		ID:        rand.Int(),
		Key:       key,
		Value:     value,
		CreatedAt: time.Now(),
	}
	writeSecret(secret)
	return secret, nil
}

func UpdateSecret(id int, key string, value string) error {
	secret, err := GetSecret(key)
	if err != nil {
		return err
	}

	secret.Key = key
	secret.Value = value
	return writeSecret(secret)
}

func DeleteSecret(id int) error {
	return deleteSecret(id)
}
