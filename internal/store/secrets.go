package store

import (
	"math/rand"
	"time"
)

type Secret struct {
	ID        int       `json:"id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

func ListSecrets() []Secret {
	return []Secret{
		// TODO - Implement from Storage
		{ID: 1, Key: "Secret 1", Value: "Value 1", CreatedAt: time.Now()},
		{ID: 2, Key: "Secret 2", Value: "Value 2", CreatedAt: time.Now()},
		{ID: 3, Key: "Secret 3", Value: "Value 3", CreatedAt: time.Now()},
	}
}

func CreateSecret(key string, value string) (Secret, error) {
	// TODO - Implement the logic to create a new secret
	return Secret{
		ID:        rand.Int(),
		Key:       key,
		Value:     value,
		CreatedAt: time.Now(),
	}, nil
}

func UpdateSecret(id int, key string, value string) error {
	// TODO - Implement the logic to update a secret by ID
	return nil
}

func DeleteSecret(id int) error {
	// TODO - Implement the logic to delete a secret by ID
	return nil
}
