package paycode

import (
	"encoding/base64"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

var testScanPayKey = []byte("wallet-demo-scanpay-key-00000001")

func TestScanPayRoundTrip(t *testing.T) {
	expires := time.Now().UTC().Add(time.Hour).Truncate(time.Second)
	want := &model.PayCodeScanPay{
		ID:      1397453561595432960,
		Amount:  decimal.RequireFromString("100000.12"),
		Expired: &expires,
	}
	s := newScanPayWithKey(testScanPayKey, "v1")

	encoded, err := s.Encode(want)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	got, err := s.Decode(encoded, &model.PayCodeScanPay{})
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("round trip got = %#v, want %#v", got, want)
	}
}

func TestScanPayRejectsExpiredCode(t *testing.T) {
	expires := time.Now().UTC().Add(-time.Minute)
	s := newScanPayWithKey(testScanPayKey, "v1")
	encoded, err := s.Encode(&model.PayCodeScanPay{ID: 1, Amount: decimal.NewFromInt(100), Expired: &expires})
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	_, err = s.Decode(encoded, &model.PayCodeScanPay{})
	if !errors.Is(err, errs.PayCodeScanPayExpired) {
		t.Fatalf("Decode() error = %v, want %v", err, errs.PayCodeScanPayExpired)
	}
}

func TestScanPayRejectsTampering(t *testing.T) {
	expires := time.Now().UTC().Add(time.Hour)
	s := newScanPayWithKey(testScanPayKey, "v1")
	encoded, err := s.Encode(&model.PayCodeScanPay{ID: 1, Amount: decimal.NewFromInt(100), Expired: &expires})
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	parts := strings.SplitN(encoded, ".", 2)
	ciphertext, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatal(err)
	}
	ciphertext[len(ciphertext)-1] ^= 1
	encoded = parts[0] + "." + base64.RawURLEncoding.EncodeToString(ciphertext)
	if _, err := s.Decode(encoded, &model.PayCodeScanPay{}); !errors.Is(err, errs.PayCodeScanPayDecodeError) {
		t.Fatalf("Decode() error = %v, want %v", err, errs.PayCodeScanPayDecodeError)
	}
}

func TestScanPayDemoFixture(t *testing.T) {
	const fixture = "v1.fecLF2eHjRwK19DTCePRfOlmcrBoXN29ZU5wzuP2f3P6IgIn6hNEUDTnC9SRkRRHs49nvfYZlJOzq-h_bUbhCudJwpRYm3z_5mj-sO9qYZjY0be5UiBO_0i3HFicd-tvhXvy5KfKtgA0mPMRgjsk6xzxOuTnw4TfOAu4P5utop9u_zLKiLVwQx0zN7PFSFqbxvVTRQXd8be1xOXO2KI"

	got, err := newScanPayWithKey(testScanPayKey, "v1").Decode(fixture, &model.PayCodeScanPay{})
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if got.ID != 900000000000000101 || !got.Amount.Equal(decimal.NewFromInt(100)) {
		t.Fatalf("decoded fixture = %#v", got)
	}
}
