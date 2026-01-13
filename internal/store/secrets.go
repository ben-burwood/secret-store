package store

import (
	"encoding/json"
	"os"
	"secret-store/internal/secret"
	"time"
)

const SecretsFile = storeDirectory + "/secrets.json"

func loadSecrets() (secret.Secrets, error) {
	if _, err := os.Stat(SecretsFile); os.IsNotExist(err) {
		// File Does Not Exist
		return secret.Secrets{}, nil
	}

	data, err := os.ReadFile(SecretsFile)
	if err != nil {
		return nil, err
	}
	secrets := secret.Secrets{}
	err = json.Unmarshal(data, &secrets)
	if err != nil {
		return nil, err
	}
	return secrets, nil
}

func ListSecrets() (secret.Secrets, error) {
	return loadSecrets()
}

func GetSecret(key string) (secret.Secret, error) {
	secrets, err := ListSecrets()
	if err != nil {
		return secret.Secret{}, err
	}
	for _, secret := range secrets {
		if secret.Key == key {
			return secret, nil
		}
	}
	return secret.Secret{}, secret.ErrSecretNotFound
}

func CreateSecret(key string, value string) (secret.Secret, error) {
	allSecrets, err := ListSecrets()
	if err != nil {
		return secret.Secret{}, err
	}
	largestSecretId := secret.Secrets(allSecrets).LargestId()

	newSecret := secret.Secret{
		ID:        largestSecretId + 1,
		Key:       key,
		Value:     value,
		CreatedAt: time.Now(),
	}
	secrets, err := loadSecrets()
	if err != nil {
		return secret.Secret{}, err
	}
	secrets = append(secrets, newSecret)
	data, err := json.Marshal(secrets)
	if err != nil {
		return secret.Secret{}, err
	}
	err = os.WriteFile(SecretsFile, data, 0644)
	if err != nil {
		return secret.Secret{}, err
	}
	return newSecret, nil
}

func UpdateSecret(id int, key string, value string) error {
	updateSecret := secret.Secret{
		ID:        id,
		Key:       key,
		Value:     value,
		CreatedAt: time.Now(),
	}

	secrets, err := loadSecrets()
	if err != nil {
		return err
	}
	for i, s := range secrets {
		if s.ID == updateSecret.ID {
			secrets[i] = updateSecret
			break
		}
	}
	data, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	return os.WriteFile(SecretsFile, data, 0644)
}

func DeleteSecret(id int) error {
	secrets, err := loadSecrets()
	if err != nil {
		return err
	}
	for i, s := range secrets {
		if s.ID == id {
			secrets = append(secrets[:i], secrets[i+1:]...)
			break
		}
	}
	data, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	return os.WriteFile(SecretsFile, data, 0644)
}

func WipeSecrets() error {
	return os.WriteFile(SecretsFile, []byte("[]"), 0644)
}
