package geetestsdk

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"
)

type cipherAES struct {
}

func (c cipherAES) EncryptString(key string, content string) (*string, error) {
	ciphertext, err := c.Encrypt([]byte(key), []byte(content))
	if err != nil {
		return nil, err
	}

	r := fmt.Sprintf("%x", ciphertext)
	return &r, nil
}

// Encrypt 使用 AES 加密
func (c cipherAES) Encrypt(key []byte, content []byte) (ciphertext []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("encrypt failed")
		}
	}()

	keyByte := c.fillKey(key)

	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()

	contentByte := c.pkcs7Padding(content, blockSize)
	ciphertext = make([]byte, len(contentByte))

	blockMode := cipher.NewCBCEncrypter(block, keyByte[:blockSize])
	blockMode.CryptBlocks(ciphertext, contentByte)

	return ciphertext, nil
}

func (c cipherAES) DecryptString(key string, ciphertext string) (*string, error) {
	ciphertextByte, err := hex.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	content, err := c.Decrypt([]byte(key), ciphertextByte)
	if err != nil {
		return nil, err
	}

	r := string(content)
	return &r, nil
}

// Decrypt 使用 AES 解密
func (c cipherAES) Decrypt(key []byte, ciphertext []byte) (context []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("decrypt failed")
		}
	}()

	keyByte := c.fillKey(key)

	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()

	origData := make([]byte, len(ciphertext))

	blockMode := cipher.NewCBCDecrypter(block, keyByte[:blockSize])
	blockMode.CryptBlocks(origData, ciphertext)

	return c.pkcs7UnPadding(origData), nil
}

func (c cipherAES) fillKey(key []byte) []byte {
	k := string(key)
	if diff := 16 - len(k); diff > 0 {
		return c.appendZeroToKey(k, diff)
	}

	if diff := 24 - len(k); diff > 0 {
		return c.appendZeroToKey(k, diff)
	}

	if diff := 32 - len(k); diff > 0 {
		return c.appendZeroToKey(k, diff)
	}

	return key
}

func (c cipherAES) appendZeroToKey(key string, count int) []byte {
	for i := 0; i < count; i++ {
		key += "z"
	}
	return []byte(key)
}

func (c cipherAES) pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}
func (c cipherAES) pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}
