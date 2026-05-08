package crypto

import (
	"errors"

	"github.com/fernet/fernet-go"
)

type Cipher struct {
	key *fernet.Key
}

func New(keyStr string) (*Cipher, error) {
	k, err := fernet.DecodeKey(keyStr)
	if err != nil {
		return nil, err
	}
	return &Cipher{key: k}, nil
}

func (c *Cipher) Encrypt(plain string) (string, error) {
	tok, err := fernet.EncryptAndSign([]byte(plain), c.key)
	if err != nil {
		return "", err
	}
	return string(tok), nil
}

// Decrypt verifies and decrypts a Fernet token. Passing ttl=0 disables the
// timestamp check, matching Python's Fernet.decrypt() default behaviour.
func (c *Cipher) Decrypt(token string) (string, error) {
	plain := fernet.VerifyAndDecrypt([]byte(token), 0, []*fernet.Key{c.key})
	if plain == nil {
		return "", errors.New("invalid fernet token")
	}
	return string(plain), nil
}
