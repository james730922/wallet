package paycode

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

var errInvalidCiphertext = errors.New("invalid ciphertext")

type cipherAES struct{}

// Encrypt uses AES-GCM and prefixes the ciphertext with a fresh random nonce.
func (cipherAES) Encrypt(key, content, additionalData []byte) ([]byte, error) {
	gcm, err := newGCM(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, content, additionalData), nil
}

// Decrypt verifies the GCM authentication tag before returning plaintext.
func (cipherAES) Decrypt(key, ciphertext, additionalData []byte) ([]byte, error) {
	gcm, err := newGCM(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < gcm.NonceSize()+gcm.Overhead() {
		return nil, errInvalidCiphertext
	}

	nonce := ciphertext[:gcm.NonceSize()]
	plaintext, err := gcm.Open(nil, nonce, ciphertext[gcm.NonceSize():], additionalData)
	if err != nil {
		return nil, errInvalidCiphertext
	}
	return plaintext, nil
}

func newGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}
