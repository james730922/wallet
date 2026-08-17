package conf

import (
	"encoding/base64"
	"errors"
	"os"
)

const (
	scanPayKeyEnv     = "WALLET_SCANPAY_KEY"
	scanPayVersionEnv = "WALLET_SCANPAY_KEY_VERSION"
)

type ScanPayCryptoHandler struct{}

// GetKey reads a base64-encoded AES-256 key from the environment or local config.
func (ScanPayCryptoHandler) GetKey() ([]byte, error) {
	encoded := os.Getenv(scanPayKeyEnv)
	if encoded == "" {
		encoded = appConf.v.GetString("scanpay_crypto.key")
	}
	key, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, errors.New("scanpay_crypto.key must be valid base64")
	}
	if len(key) != 32 {
		return nil, errors.New("scanpay_crypto.key must decode to 32 bytes")
	}
	return key, nil
}

func (ScanPayCryptoHandler) GetVersion() string {
	if version := os.Getenv(scanPayVersionEnv); version != "" {
		return version
	}
	return appConf.v.GetString("scanpay_crypto.version")
}
