package secret

import (
	"errors"
	"time"
)

var ErrSecretNotFound = errors.New("secret not found")

type Secret struct {
	ID        int       `json:"id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

type Secrets []Secret

func (s Secrets) LargestId() int {
	maxId := 0
	for _, secret := range s {
		if secret.ID > maxId {
			maxId = secret.ID
		}
	}
	return maxId
}

func (s Secrets) GetByID(id int) (Secret, error) {
	for _, secret := range s {
		if secret.ID == id {
			return secret, nil
		}
	}
	return Secret{}, ErrSecretNotFound
}
