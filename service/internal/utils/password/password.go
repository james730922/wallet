package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/crypto/argon2"

	"github.com/james730922/wallet/service/internal/utils/signs"
)

const (
	memory      = 64 * 1024
	iterations  = 3
	parallelism = 1
	saltLength  = 16
	keyLength   = 32
)

var errInvalidHash = errors.New("invalid password hash")

var argon2Limiter = struct {
	sync.RWMutex
	slots chan struct{}
}{
	slots: make(chan struct{}, 4),
}

// SetMaxConcurrentArgon2 configures the process-wide memory bound for Argon2.
// It is expected to be called during application startup, before requests run.
func SetMaxConcurrentArgon2(max int) {
	if max <= 0 {
		max = 4
	}
	argon2Limiter.Lock()
	argon2Limiter.slots = make(chan struct{}, max)
	argon2Limiter.Unlock()
}

func Hash(plainText string) (string, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2IDKey([]byte(plainText), salt, iterations, memory, parallelism, keyLength)
	return encode(salt, hash), nil
}

// Verify accepts Argon2id hashes and the legacy SHA-512 format. needsUpgrade is
// true only when a valid legacy hash should be replaced after authentication.
func Verify(encoded, plainText, legacySalt string) (valid, needsUpgrade bool, err error) {
	if !strings.HasPrefix(encoded, "$argon2id$") {
		legacy := signs.Hex(plainText, legacySalt)
		valid := subtle.ConstantTimeCompare([]byte(encoded), []byte(legacy)) == 1
		return valid, valid, nil
	}

	params, salt, expected, err := decode(encoded)
	if err != nil {
		return false, false, err
	}
	actual := argon2IDKey([]byte(plainText), salt, params.iterations, params.memory, params.parallelism, uint32(len(expected)))
	return subtle.ConstantTimeCompare(actual, expected) == 1, false, nil
}

func argon2IDKey(password, salt []byte, iterations, memory uint32, parallelism uint8, keyLength uint32) []byte {
	argon2Limiter.RLock()
	slots := argon2Limiter.slots
	argon2Limiter.RUnlock()

	slots <- struct{}{}
	defer func() { <-slots }()
	return argon2.IDKey(password, salt, iterations, memory, parallelism, keyLength)
}

type parameters struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
}

func encode(salt, hash []byte) string {
	b64 := base64.RawStdEncoding
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		memory,
		iterations,
		parallelism,
		b64.EncodeToString(salt),
		b64.EncodeToString(hash),
	)
}

func decode(encoded string) (parameters, []byte, []byte, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return parameters{}, nil, nil, errInvalidHash
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil || version != argon2.Version {
		return parameters{}, nil, nil, errInvalidHash
	}

	params := parameters{}
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &params.memory, &params.iterations, &params.parallelism); err != nil {
		return parameters{}, nil, nil, errInvalidHash
	}
	if params.memory < 8*1024 || params.memory > 256*1024 || params.iterations < 1 || params.iterations > 10 || params.parallelism < 1 || params.parallelism > 8 {
		return parameters{}, nil, nil, errInvalidHash
	}

	b64 := base64.RawStdEncoding
	salt, err := b64.DecodeString(parts[4])
	if err != nil || len(salt) < 8 || len(salt) > 64 {
		return parameters{}, nil, nil, errInvalidHash
	}
	hash, err := b64.DecodeString(parts[5])
	if err != nil || len(hash) < 16 || len(hash) > 64 {
		return parameters{}, nil, nil, errInvalidHash
	}

	return params, salt, hash, nil
}
