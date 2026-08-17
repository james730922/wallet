package eventbus

import "sync"

func newReceiveRepository() IReceiveRepository {
	return &receiveRepository{
		receivers: make(map[string]IReceiver),
	}
}

type IReceiveRepository interface {
	Store(topic string, receiver IReceiver)
	Load(topic string) (IReceiver, bool)
	Delete(topic string)
}

type receiveRepository struct {
	receivers map[string]IReceiver
	mx        sync.RWMutex
}

func (rr *receiveRepository) Store(topic string, receiver IReceiver) {
	rr.mx.Lock()
	rr.receivers[topic] = receiver
	rr.mx.Unlock()
}

func (rr *receiveRepository) Load(topic string) (IReceiver, bool) {
	rr.mx.RLock()
	receive, ok := rr.receivers[topic]
	rr.mx.RUnlock()
	return receive, ok
}

func (rr *receiveRepository) Delete(topic string) {
	rr.mx.Lock()
	delete(rr.receivers, topic)
	rr.mx.Unlock()
}
