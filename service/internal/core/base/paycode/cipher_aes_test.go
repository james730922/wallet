package paycode

import (
	"bytes"
	"testing"
)

func TestAESCipherRoundTrip(t *testing.T) {
	key := []byte("wallet-demo-scanpay-key-00000001")
	plaintext := []byte("matt")
	aad := []byte("wallet-scanpay:v1")
	ciphertext, err := (cipherAES{}).Encrypt(key, plaintext, aad)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}
	got, err := (cipherAES{}).Decrypt(key, ciphertext, aad)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}
	if !bytes.Equal(got, plaintext) {
		t.Fatalf("Decrypt() got = %q, want %q", got, plaintext)
	}
}

func TestAESCipherRejectsTampering(t *testing.T) {
	key := []byte("wallet-demo-scanpay-key-00000001")
	aad := []byte("wallet-scanpay:v1")
	ciphertext, err := (cipherAES{}).Encrypt(key, []byte("matt"), aad)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}
	ciphertext[len(ciphertext)-1] ^= 1
	if _, err := (cipherAES{}).Decrypt(key, ciphertext, aad); err == nil {
		t.Fatal("Decrypt() accepted tampered ciphertext")
	}
}
