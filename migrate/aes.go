package migrate

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// indicates key size is too small.
var errKeySize = errors.New("encryption key must be 32 bytes")

// helper function parses the encryption key.
func parseKey(key string) (cipher.Block, error) {
	if len(key) != 32 {
		return nil, errKeySize
	}
	b := []byte(key)
	return aes.NewCipher(b)
}

// helper function to encrypt secrets.
func encrypt(block cipher.Block, plaintext string) ([]byte, error) {
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, []byte(plaintext), nil), nil
}
