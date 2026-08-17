package scanpaymember

import (
	"bytes"
	"encoding/base64"
	"errors"
	"testing"

	"github.com/james730922/wallet/service/internal/utils/errs"
)

func TestDecodeQRCodeImageSizeLimit(t *testing.T) {
	atLimit := bytes.Repeat([]byte{1}, int(maxQRCodeImageBytes))
	decoded, err := decodeQRCodeImage(base64.StdEncoding.EncodeToString(atLimit))
	if err != nil {
		t.Fatalf("image at limit rejected: %v", err)
	}
	if len(decoded) != len(atLimit) {
		t.Fatalf("decoded size = %d, want %d", len(decoded), len(atLimit))
	}

	overLimit := bytes.Repeat([]byte{1}, int(maxQRCodeImageBytes)+1)
	_, err = decodeQRCodeImage(base64.StdEncoding.EncodeToString(overLimit))
	if !errors.Is(err, errs.ScanQRCodeImageTooLarge) {
		t.Fatalf("oversized image error = %v, want %v", err, errs.ScanQRCodeImageTooLarge)
	}
}

func TestDecodeQRCodeImageRejectsInvalidBase64(t *testing.T) {
	_, err := decodeQRCodeImage("not-valid-base64")
	if !errors.Is(err, errs.ScanQRCodeVerifyFailed) {
		t.Fatalf("invalid Base64 error = %v, want %v", err, errs.ScanQRCodeVerifyFailed)
	}
}
