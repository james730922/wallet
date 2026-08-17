package password

import (
	"testing"

	"github.com/james730922/wallet/service/internal/utils/conf"
	"github.com/james730922/wallet/service/internal/utils/signs"
)

func TestHashAndVerify(t *testing.T) {
	encoded, err := Hash("demo1234")
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}

	valid, needsUpgrade, err := Verify(encoded, "demo1234", "")
	if err != nil || !valid || needsUpgrade {
		t.Fatalf("Verify() = (%v, %v, %v), want (true, false, nil)", valid, needsUpgrade, err)
	}

	valid, _, err = Verify(encoded, "wrong-password", "")
	if err != nil || valid {
		t.Fatalf("Verify(wrong) = (%v, %v), want (false, nil)", valid, err)
	}
}

func TestVerifyLegacyHash(t *testing.T) {
	conf.Mock()
	const salt = "legacy-salt"
	legacy := signs.Hex("demo1234", salt)
	valid, needsUpgrade, err := Verify(legacy, "demo1234", salt)
	if err != nil || !valid || !needsUpgrade {
		t.Fatalf("Verify(legacy) = (%v, %v, %v), want (true, true, nil)", valid, needsUpgrade, err)
	}
}

func TestSetMaxConcurrentArgon2(t *testing.T) {
	defer SetMaxConcurrentArgon2(4)
	SetMaxConcurrentArgon2(2)

	argon2Limiter.RLock()
	capacity := cap(argon2Limiter.slots)
	argon2Limiter.RUnlock()
	if capacity != 2 {
		t.Fatalf("Argon2 limiter capacity = %d, want 2", capacity)
	}
}
