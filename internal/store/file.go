package store

import (
	"encoding/json"
	"os"
)

const SecretsFile = "secrets.json"

func loadSecrets() (Secrets, error) {
	if _, err := os.Stat(SecretsFile); os.IsNotExist(err) {
		// File Does Not Exist
		return Secrets{}, nil
	}

	data, err := os.ReadFile(SecretsFile)
	if err != nil {
		return nil, err
	}
	secrets := Secrets{}
	err = json.Unmarshal(data, &secrets)
	if err != nil {
		return nil, err
	}
	return secrets, nil
}

func writeSecret(secret Secret) error {
	secrets, err := loadSecrets()
	if err != nil {
		return err
	}
	secrets = append(secrets, secret)
	data, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	return os.WriteFile(SecretsFile, data, 0644)
}

func updateSecret(secret Secret) error {
	secrets, err := loadSecrets()
	if err != nil {
		return err
	}
	for i, s := range secrets {
		if s.ID == secret.ID {
			secrets[i] = secret
			break
		}
	}
	data, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	return os.WriteFile(SecretsFile, data, 0644)
}

func deleteSecret(id int) error {
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
