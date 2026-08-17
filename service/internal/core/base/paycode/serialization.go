package paycode

import (
	"bytes"
	"encoding/gob"
)

type serialization struct {
}

// Encode 物件轉 []byte
func (s *serialization) Encode(content interface{}) ([]byte, error) {
	buf := bytes.Buffer{}

	err := gob.NewEncoder(&buf).Encode(content)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Decode []byte 轉物件
func (s *serialization) Decode(content []byte, t interface{}) error {
	buf := bytes.Buffer{}

	buf.Write(content)

	err := gob.NewDecoder(&buf).Decode(t)
	if err != nil {
		return err
	}

	return nil
}
