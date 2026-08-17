package tools

import (
	"bytes"
	"encoding/json"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

func JsonMarshal(obj interface{}) ([]byte, error) {
	return jsonMarshal(obj)
}

func JsonPBMarshal(obj proto.Message) ([]byte, error) {
	return jsonPBMarshal(obj)
}

func JsonMarshalString(obj interface{}) string {
	buf, err := jsonMarshal(obj)
	if err != nil {
		return err.Error()
	}

	return string(buf)
}

func JsonPBMarshalString(obj proto.Message) string {
	buf, err := jsonPBMarshal(obj)
	if err != nil {
		return err.Error()
	}

	return string(buf)
}

func jsonMarshal(obj interface{}) ([]byte, error) {
	if o, ok := obj.(proto.Message); ok {
		return jsonPBMarshal(o)
	}

	buf, err := json.Marshal(obj)
	if err != nil {
		return nil, errors.Errorf("err: %s", err)
	}

	return buf, nil
}

func jsonPBMarshal(obj proto.Message) ([]byte, error) {
	m := jsonpb.Marshaler{
		EmitDefaults: true,
		EnumsAsInts:  true,
		OrigName:     true,
	}

	var buf bytes.Buffer
	if err := m.Marshal(&buf, obj); err != nil {
		return nil, errors.Errorf("err: %s", err)
	}

	return buf.Bytes(), nil
}
