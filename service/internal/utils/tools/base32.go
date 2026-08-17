package tools

import (
	"encoding/base32"
	"encoding/binary"
	"errors"
)

var Base32WithInt64 = newBase32WithInt64()

func newBase32WithInt64() *base32WithInt64 {
	return &base32WithInt64{
		base32Std: base32.StdEncoding.WithPadding(-1),
	}
}

type base32WithInt64 struct {
	base32Std *base32.Encoding
}

func (b base32WithInt64) Encode(memberID int64) string {
	return b.base32Std.EncodeToString(b.int64ToBytes(memberID))
}

func (b base32WithInt64) Decode(memberID string) (int64, error) {
	buf, err := b.base32Std.DecodeString(memberID)
	if err != nil {
		return 0, err
	}

	if len(buf) != 8 {
		return 0, errors.New("not type int64, len([]byte) != 8")
	}

	return b.bytesToInt64(buf), nil
}

func (b base32WithInt64) int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(i))
	return buf
}

func (b base32WithInt64) bytesToInt64(buf []byte) int64 {
	return int64(binary.LittleEndian.Uint64(buf))
}
