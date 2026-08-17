package paycode

import (
	"reflect"
	"testing"
)

func TestSerializationStringRoundTrip(t *testing.T) {
	s := serialization{}
	want := "demo"
	encoded, err := s.Encode(want)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}

	var got string
	if err := s.Decode(encoded, &got); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if got != want {
		t.Fatalf("round trip got = %q, want %q", got, want)
	}
}

func TestSerializationStructRoundTrip(t *testing.T) {
	type content struct {
		ID     string
		Amount float64
	}
	s := serialization{}
	want := content{ID: "demo-order", Amount: 7788.8877}
	encoded, err := s.Encode(want)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}

	var got content
	if err := s.Decode(encoded, &got); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("round trip got = %#v, want %#v", got, want)
	}
}
