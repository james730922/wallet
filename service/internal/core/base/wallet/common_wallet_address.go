package wallet

import (
	"encoding/binary"
	"fmt"

	"github.com/google/uuid"
)

var salt = []byte("2aija8cd6f7hijklmnop")

func XorBytes(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("length of byte slices is not equivalent: %d != %d", len(a), len(b))
	}

	buf := make([]byte, len(a))

	for i, _ := range a {
		buf[i] = a[i] ^ b[i]
	}

	return buf, nil
}

func GenerateWalletAddress(i int) string {
	id := uuid.New()
	buf, _ := id.MarshalBinary()

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(i))

	buf = append(buf, b...)

	tmp, _ := XorBytes(buf, salt)
	return fmt.Sprintf("0x%x", tmp)
}
