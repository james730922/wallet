package paycode

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/conf"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

const scanPayAADPrefix = "wallet-scanpay:"

func newScanPay() *scanPay {
	key, err := conf.ScanPayCrypto().GetKey()
	if err != nil {
		panic(err)
	}
	return newScanPayWithKey(key, conf.ScanPayCrypto().GetVersion())
}

func newScanPayWithKey(key []byte, version string) *scanPay {
	return &scanPay{aesKey: append([]byte(nil), key...), version: version}
}

type IScanPay interface {
	Encode(scanPay *model.PayCodeScanPay) (string, error)
	Decode(packetText string, scanPay *model.PayCodeScanPay) (*model.PayCodeScanPay, error)
}

type scanPay struct {
	cipher        cipherAES
	serialization serialization
	aesKey        []byte
	version       string
}

func (s scanPay) Encode(scanPay *model.PayCodeScanPay) (string, error) {
	payload, err := s.serialization.Encode(scanPay)
	if err != nil {
		logger.ApLog().Warnf("scanpay encode serialization failed: %s", err)
		return "", errs.PayCodeScanPayEncodeError
	}

	ciphertext, err := s.cipher.Encrypt(s.aesKey, payload, s.additionalData())
	if err != nil {
		logger.ApLog().Warnf("scanpay encryption failed: %s", err)
		return "", errs.PayCodeScanPayEncodeError
	}

	return s.version + "." + base64.RawURLEncoding.EncodeToString(ciphertext), nil
}

func (s scanPay) Decode(packetText string, scanPay *model.PayCodeScanPay) (*model.PayCodeScanPay, error) {
	parts := strings.SplitN(packetText, ".", 2)
	if len(parts) != 2 || parts[0] != s.version || parts[1] == "" {
		return nil, errs.PayCodeScanPayDecodeError
	}

	ciphertext, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errs.PayCodeScanPayDecodeError
	}
	content, err := s.cipher.Decrypt(s.aesKey, ciphertext, s.additionalData())
	if err != nil {
		return nil, errs.PayCodeScanPayDecodeError
	}

	if err := s.serialization.Decode(content, scanPay); err != nil {
		logger.ApLog().Debugf("scanpay decode serialization failed: %s", err)
		return nil, errs.PayCodeScanPayDecodeError
	}
	if scanPay.Expired != nil && scanPay.Expired.Before(time.Now().UTC()) {
		return nil, errs.PayCodeScanPayExpired
	}

	return scanPay, nil
}

func (s scanPay) additionalData() []byte {
	return []byte(scanPayAADPrefix + s.version)
}
