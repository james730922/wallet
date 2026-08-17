package auth

import (
	"errors"
	"strings"
	"testing"

	"github.com/james730922/wallet/service/internal/utils/errs"
)

func TestLoginGuardAttemptThresholds(t *testing.T) {
	guard := &loginGuard{accountMax: 5, ipMax: 30, captchaAfter: 3}

	requireCaptcha, err := guard.evaluateAttempts(2, 10)
	if err != nil || requireCaptcha {
		t.Fatalf("below threshold: captcha=%t err=%v", requireCaptcha, err)
	}
	requireCaptcha, err = guard.evaluateAttempts(3, 10)
	if err != nil || !requireCaptcha {
		t.Fatalf("captcha threshold: captcha=%t err=%v", requireCaptcha, err)
	}
	if _, err := guard.evaluateAttempts(5, 10); !errors.Is(err, errs.CommonFrequentOperationError) {
		t.Fatalf("account limit error = %v", err)
	}
	if _, err := guard.evaluateAttempts(1, 30); !errors.Is(err, errs.CommonFrequentOperationError) {
		t.Fatalf("IP limit error = %v", err)
	}
}

func TestLoginAttemptKeysDoNotContainPII(t *testing.T) {
	accountKey, ipKey := loginAttemptKeys("86", "13800000000", "192.0.2.1")
	for _, key := range []string{accountKey, ipKey} {
		if strings.Contains(key, "13800000000") || strings.Contains(key, "192.0.2.1") {
			t.Fatalf("rate-limit key leaks PII: %s", key)
		}
	}
}

func TestRegistrationAttemptKeysDoNotContainRawIdentifiers(t *testing.T) {
	ipKey, deviceKey := registrationAttemptKeys("192.0.2.1", "device-secret")
	for _, key := range []string{ipKey, deviceKey} {
		if strings.Contains(key, "192.0.2.1") || strings.Contains(key, "device-secret") {
			t.Fatalf("registration rate-limit key leaks identifier: %s", key)
		}
	}
}

func TestRegistrationDeviceKeyIsScopedByIP(t *testing.T) {
	_, first := registrationAttemptKeys("192.0.2.1", "same-device")
	_, second := registrationAttemptKeys("192.0.2.2", "same-device")
	if first == second {
		t.Fatal("device keys from different IPs must not collide")
	}
}

func TestPasswordChangeOwnership(t *testing.T) {
	if err := validatePasswordChangeOwnership(7, 7); err != nil {
		t.Fatalf("owner rejected: %v", err)
	}
	if err := validatePasswordChangeOwnership(7, 8); !errors.Is(err, errs.AuthOperationForbidden) {
		t.Fatalf("cross-member update error = %v", err)
	}
}
