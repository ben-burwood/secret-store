package store

import (
	"os"
)

const storeDirectory = "store"

func InitialiseStore() error {
	// Write Dir if not Exist
	if _, err := os.Stat(storeDirectory); os.IsNotExist(err) {
		if err := os.Mkdir(storeDirectory, 0755); err != nil {
			return err
		}
	}
	return nil
}
