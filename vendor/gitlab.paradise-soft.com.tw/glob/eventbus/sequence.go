package eventbus

import "sync"

func newSequence() ISequence {
	return &sequence{}
}

type ISequence interface {
	Next() EntryID
}

type sequence struct {
	now EntryID
	mx  sync.Mutex
}

func (seq *sequence) Next() EntryID {
	seq.mx.Lock()
	defer seq.mx.Unlock()
	seq.now++
	return seq.now
}
