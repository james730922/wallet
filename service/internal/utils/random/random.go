package random

import (
	"math/rand"
	"sync"
	"time"
)

var r = New(time.Now().UTC().Unix())

// Random is a random generator with convenient methods
type Random struct {
	mx        sync.Mutex
	generator *rand.Rand
}

// New returns a new Random object initialized with the given seed
func New(seed int64) *Random {
	r := &Random{}
	source := rand.NewSource(seed)
	r.generator = rand.New(source)
	return r
}

// Intn returns a random int with in [0..n]
func Intn(n int) int { return r.Intn(n) }
func (r *Random) Intn(n int) int {
	if n == 0 {
		return 0
	}

	r.mx.Lock()
	i := r.generator.Intn(n)
	r.mx.Unlock()
	return i
}

// Int32n returns a random int32 with in [0..n]
func Int32n(n int32) int32 { return r.Int32n(n) }
func (r *Random) Int32n(n int32) int32 {
	if n == 0 {
		return 0
	}

	r.mx.Lock()
	i := r.generator.Int31n(n)
	r.mx.Unlock()
	return i
}

// Int64n returns a random int64 with in [0..n]
func Int64n(n int64) int64 { return r.Int64n(n) }
func (r *Random) Int64n(n int64) int64 {
	if n == 0 {
		return 0
	}

	r.mx.Lock()
	i := r.generator.Int63n(n)
	r.mx.Unlock()
	return i
}

// Selection 打亂dictionary ，依據size 從dictionary 挑出字元並回傳
func Selection(size int, dictionary string) string { return r.Selection(size, dictionary) }
func (r *Random) Selection(size int, dictionary string) string {
	if size == 0 || len(dictionary) == 0 {
		return ""
	}

	var bytes = make([]byte, size)
	r.generator.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}

	return string(bytes)
}
