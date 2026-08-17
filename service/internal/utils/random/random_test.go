package random

import (
	"fmt"
	"testing"
	"time"
)

func TestSomething(t *testing.T) {
	ms := make(map[int64]int64)
	max := int64(10)

	r := New(time.Now().UTC().Unix())
	for i := 0; i < 10000; i++ {
		v := r.Int64n(max)
		ms[v] = ms[v] + 1
	}

	fmt.Println("ms len:", len(ms))
	for i := int64(0); i < max; i++ {
		if v, ok := ms[i]; ok {
			fmt.Printf("key: %d, val: %d\n", i, v)
		} else {

		}
	}

	for k, v := range ms {
		fmt.Printf("key: %d, val: %d\n", k, v)
	}
}
